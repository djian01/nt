package ntPinger

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// helper: basic auth header builder for proxy connect
func basicAuthFromURLUser(u *url.Userinfo) string {
	if u == nil {
		return ""
	}
	username := u.Username()
	password, _ := u.Password()
	token := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return "Basic " + token
}

// helper: is this a proxy CONNECT/tunnel error (Go's transport message)
func isProxyConnectError(err error) bool {
	if err == nil {
		return false
	}
	e := strings.ToLower(err.Error())
	// typical patterns: "proxyconnect", "bad response from proxy", "proxy error"
	if strings.Contains(e, "proxyconnect") || strings.Contains(e, "bad response from proxy") ||
		strings.Contains(e, "received http code") && strings.Contains(e, "from proxy") {
		return true
	}
	return false
}

// helper: inspect headers/body/tls to decide if response is from a proxy
func isProxyResponse(resp *http.Response, usedProxy bool, bodySnippet string) bool {
	if resp == nil || !usedProxy {
		return false
	}

	// Look at obvious headers
	h := resp.Header
	server := strings.ToLower(h.Get("Server"))
	if strings.Contains(server, "squid") || strings.Contains(server, "zscaler") ||
		strings.Contains(server, "proxy") || strings.Contains(server, "bluecoat") ||
		strings.Contains(server, "barracuda") {
		return true
	}

	// Proxy-specific headers
	if h.Get("X-Squid-Error") != "" || h.Get("Proxy-Agent") != "" ||
		strings.Contains(strings.ToLower(h.Get("Via")), "zscaler") ||
		strings.Contains(strings.ToLower(h.Get("Via")), "squid") {
		return true
	}

	// Body sniffing (small snippet) for vendor keywords
	b := strings.ToLower(bodySnippet)
	if b != "" {
		if strings.Contains(b, "zscaler") || strings.Contains(b, "x-squid-error") ||
			strings.Contains(b, "this request was blocked by") || strings.Contains(b, "access denied by proxy") {
			return true
		}
	}

	// TLS issuer check (if TLS present)
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		iss := resp.TLS.PeerCertificates[0].Issuer
		issStr := strings.ToLower(strings.Join(append(iss.Organization, iss.CommonName), " "))
		if strings.Contains(issStr, "zscaler") || strings.Contains(issStr, "zscaler, inc") ||
			strings.Contains(issStr, "bluecoat") || strings.Contains(issStr, "symantec") {
			return true
		}
	}

	return false
}

// HttpProbing performs HTTP probing with proxy-detection enhancement
func HttpProbing_1(
	Seq int,
	destHost string,
	destPort int,
	Http_path string,
	Http_scheme string,
	Http_method string,
	Http_statusCodes []HttpStatusCode,
	timeout int,
	proxyStr string,
) (PacketHTTP, error) {

	// Initial PacketHTTP
	pkt := PacketHTTP{
		Type:        "http",
		Status:      false,
		Seq:         Seq,
		DestHost:    destHost,
		DestPort:    destPort,
		Http_scheme: Http_scheme,
		Http_method: Http_method,
		Http_path:   Http_path,
	}

	// Construct the URL for HTTP or HTTPS
	testUrl := ConstructURL(Http_scheme, destHost, Http_path, destPort)

	// Create a custom Transport to ignore certificate errors
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Good defaults for probes:
		DisableCompression:  true,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// Configure proxy
	usedProxy := false
	if proxyStr != "none" && strings.TrimSpace(proxyStr) != "" {
		usedProxy = true
		u, err := url.Parse(proxyStr)
		if err != nil {
			return pkt, fmt.Errorf("invalid proxy URL: %w", err)
		}

		switch u.Scheme {
		case "http", "https":
			tr.Proxy = http.ProxyURL(u)

			// If credentials present, ensure CONNECT gets Proxy-Authorization
			if u.User != nil {
				if h := basicAuthFromURLUser(u.User); h != "" {
					if tr.ProxyConnectHeader == nil {
						tr.ProxyConnectHeader = make(http.Header, 1)
					}
					tr.ProxyConnectHeader.Set("Proxy-Authorization", h)
				}
			}

		case "socks5", "socks5h":
			// Build a SOCKS5 dialer, optionally with auth
			var auth *proxy.Auth
			if u.User != nil {
				pass, _ := u.User.Password()
				auth = &proxy.Auth{User: u.User.Username(), Password: pass}
			}

			// host:port
			addr := u.Host
			if !strings.Contains(addr, ":") {
				// default SOCKS5 port
				addr = net.JoinHostPort(addr, "1080")
			}

			d, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
			if err != nil {
				return pkt, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
			}

			// Adapt to DialContext
			tr.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				// socks5 dialer only exposes Dial, so wrap it
				return d.Dial(network, address)
			}

		default:
			return pkt, fmt.Errorf("unsupported proxy scheme %q (use http(s) or socks5)", u.Scheme)
		}
	} else {
		// If not explicitly set, honor environment proxies
		tr.Proxy = http.ProxyFromEnvironment
	}

	// Create a new HTTP client with the custom Transport and timeout
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}

	// Record the start time
	pkt.SendTime = time.Now()

	// Create a new request based on the specified HTTP method
	req, err := http.NewRequest(Http_method, testUrl, nil)
	if err != nil {
		return pkt, fmt.Errorf("failed to create request: %w", err)
	}

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		// transport-level errors (DNS, dial, proxyconnect, TLS, timeout, etc.)
		switch {
		case strings.Contains(strings.ToLower(err.Error()), "context deadline exceeded"):
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil

		case strings.Contains(strings.ToLower(err.Error()), "timeout"):
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil

		case strings.Contains(strings.ToLower(err.Error()), "handshake timeout"):
			pkt.AdditionalInfo = "Handshake_Timeout"
			return pkt, nil

		case strings.Contains(strings.ToLower(err.Error()), "network is unreachable"):
			pkt.AdditionalInfo = "Network_Unreachable"
			return pkt, nil

		case strings.Contains(strings.ToLower(err.Error()), "connection refused"):
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "Conn_Refused"
			return pkt, nil

		case strings.Contains(strings.ToLower(err.Error()), "connection reset by peer"):
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "Conn_Reset_by_Peer"
			return pkt, nil

		case strings.Contains(strings.ToLower(err.Error()), "eof"):
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "EOF"
			return pkt, nil

		default:
			// If this looks like a proxy CONNECT refusal, classify as Proxy_Block
			if isProxyConnectError(err) {
				pkt.RTT = time.Since(pkt.SendTime)
				pkt.AdditionalInfo = "Proxy_Block"
				return pkt, nil
			}
			return pkt, fmt.Errorf("failed to send request: %w", err)
		}
	}
	// Ensure body closed
	defer resp.Body.Close()

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Fill in the HTTP response code & response phrase
	pkt.Http_response_code = resp.StatusCode
	if parts := strings.SplitN(resp.Status, " ", 2); len(parts) == 2 {
		pkt.Http_response = parts[1]
	} else {
		pkt.Http_response = resp.Status
	}

	// Read a small body snippet for detection of block pages (limit 4KB)
	var bodySnippet string
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		bodySnippet = string(b)
		// reattach a no-op body if other code expects to read it later
		resp.Body = io.NopCloser(bytes.NewReader(b))
	}

	// If this is a likely block status (403, 407, 451, 401, maybe 302 redirect to portal)
	if resp.StatusCode == 403 || resp.StatusCode == 407 || resp.StatusCode == 451 || resp.StatusCode == 401 || resp.StatusCode == 302 {
		if isProxyResponse(resp, usedProxy, bodySnippet) {
			pkt.AdditionalInfo = "Proxy_Block"
		} else {
			// origin block
			pkt.AdditionalInfo = "Forbidden"
		}
	}

	// Status decision based on allowed ranges
	pkt.Status = statusAllowed(resp.StatusCode, Http_statusCodes)
	if !pkt.Status {
		if pkt.AdditionalInfo == "" {
			if desc, ok := HTTPStatusDescription[resp.StatusCode]; ok {
				pkt.AdditionalInfo = desc
			} else {
				pkt.AdditionalInfo = "StatusNotAllowed"
			}
		}
	}

	// Check for latency
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}
