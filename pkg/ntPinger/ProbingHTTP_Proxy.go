package ntPinger

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

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
	if strings.Contains(e, "forbidden") || strings.Contains(e, "proxy") {
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

	return false
}

// configureProxy updates the provided http.Transport (tr) based on the given proxy string.
// Supported proxy schemes: http, https, socks5, socks5h
// Pass "none" or an empty string to skip proxy configuration.
// proxyStr is already check before calling this function. proxyStr is NOT empty in the input.
func configureProxy(proxyStr string, tr *http.Transport) error {

	// Parse the proxy URL
	u, err := url.Parse(proxyStr)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	switch u.Scheme {
	case "http", "https":
		tr.Proxy = http.ProxyURL(u)

		// Handle HTTP(S) proxy authentication
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
			addr = net.JoinHostPort(addr, "1080")
		}

		d, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
		if err != nil {
			return fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
		}

		// Wrap DialContext to use SOCKS5 dialer
		tr.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			// socks5 dialer only exposes Dial, so wrap it
			return d.Dial(network, address)
		}

	default:
		return fmt.Errorf("unsupported proxy scheme %q (use http(s) or socks5)", u.Scheme)
	}

	return nil
}
