package ntPinger

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/proxy" // for SOCKS5
)

var HTTPStatusDescription = map[int]string{
	// 1xx — Informational
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",

	// 2xx — Success
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",

	// 3xx — Redirection
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy (Deprecated)",
	307: "Temporary Redirect",
	308: "Permanent Redirect",

	// 4xx — Client Error
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",

	// 5xx — Server Error
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

// func: httpProbingRun
func httpProbingRun(p *Pinger, errChan chan<- error) {

	// Create a channel to listen for SIGINT (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	defer close(interruptChan)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Sequence number
	Seq := 0

	// forLoopEnds Flag
	forLoopEnds := false

	// count
	if p.InputVars.Count == 0 { // no specificed test count
		for {
			// Loop End Signal
			if forLoopEnds {
				break
			}

			// Pinger end Singal
			if p.PingerEnd {
				interruptChan <- os.Interrupt //send interrupt to interruptChan
			}

			// Perform HTTP probing
			pkt, err := HttpProbing(
				Seq,
				p.InputVars.DestHost,
				p.InputVars.DestPort,
				p.InputVars.Http_path,
				p.InputVars.Http_scheme,
				p.InputVars.Http_method,
				p.InputVars.Http_statusCodes,
				p.InputVars.Timeout,
				p.InputVars.Http_proxy,
			)
			if err != nil {
				errChan <- err
			}

			// Update statistics and send packet
			p.UpdateStatistics(&pkt)
			pkt.UpdateStatistics(p.Stat)
			p.ProbeChan <- &pkt
			Seq++

			// sleep for interval
			select {
			case <-time.After(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT)):
				// wait for SleepTime
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				forLoopEnds = true
				close(p.ProbeChan)
			}
		}

	} else { // specificed test count
		for i := 0; i < p.InputVars.Count; i++ {

			if forLoopEnds {
				break
			}

			// Perform HTTP probing
			pkt, err := HttpProbing(
				Seq,
				p.InputVars.DestHost,
				p.InputVars.DestPort,
				p.InputVars.Http_path,
				p.InputVars.Http_scheme,
				p.InputVars.Http_method,
				p.InputVars.Http_statusCodes,
				p.InputVars.Timeout,
				p.InputVars.Http_proxy,
			)
			if err != nil {
				errChan <- err
			}

			// Update statistics and send packet
			p.UpdateStatistics(&pkt)
			pkt.UpdateStatistics(p.Stat)
			p.ProbeChan <- &pkt
			Seq++

			// check the last loop of the probing, close probeChan
			if i == (p.InputVars.Count - 1) {
				close(p.ProbeChan)
			}

			// sleep for interval
			select {
			case <-time.After(GetSleepTime(pkt.Status, p.InputVars.Interval, pkt.RTT)):
				// wait for SleepTime
			case <-interruptChan: // case interruptChan, close the channel & break the loop
				forLoopEnds = true
				close(p.ProbeChan)
			}
		}
	}
}

// HttpProbing performs HTTP probing with the ability to choose HTTP methods and ignore certificate errors
func HttpProbing(
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
	if proxyStr != "none" {
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
	if err != nil { // This happens before or instead of a valid HTTP response — Network or Protocol Failure
		switch {
		case strings.Contains(err.Error(), "context deadline exceeded"):
			// Timeout error
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil

		case strings.Contains(err.Error(), "timeout"):
			// Timeout error
			pkt.AdditionalInfo = "Conn_Timeout"
			return pkt, nil

		case strings.Contains(err.Error(), "handshake timeout"):
			// Connection handshake timeout
			pkt.AdditionalInfo = "Handshake_Timeout"
			return pkt, nil

		case strings.Contains(err.Error(), "network is unreachable"):
			// Connection network is unreachable
			pkt.AdditionalInfo = "Network_Unreachable"
			return pkt, nil

		case strings.Contains(err.Error(), "connection refused"):
			// Connection refused error
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "Conn_Refused"
			return pkt, nil

		case strings.Contains(err.Error(), "connection reset by peer"):
			// Connection reset by peer
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "Conn_Reset_by_Peer"
			return pkt, nil

		case strings.Contains(err.Error(), "EOF"):
			// EOF
			pkt.RTT = time.Since(pkt.SendTime)
			pkt.AdditionalInfo = "EOF"
			return pkt, nil

		default:
			fmt.Println(err)
			return pkt, fmt.Errorf("failed to send request: %w", err)
		}
	}
	defer resp.Body.Close()

	// Calculate RTT
	pkt.RTT = time.Since(pkt.SendTime)

	// Fill in the HTTP response code & response phase
	pkt.Http_response_code = resp.StatusCode
	if parts := strings.SplitN(resp.Status, " ", 2); len(parts) == 2 {
		pkt.Http_response = parts[1]
	} else {
		pkt.Http_response = resp.Status
	}

	// Status decision based on allowed ranges
	pkt.Status = statusAllowed(resp.StatusCode, Http_statusCodes)
	if !pkt.Status {
		if pkt.AdditionalInfo == "" {
			// Additional Info for Failures
			if desc, ok := HTTPStatusDescription[resp.StatusCode]; ok {
				pkt.AdditionalInfo = desc
			} else {
				pkt.AdditionalInfo = "StatusNotAllowed"
			}
		}
	}

	// Check for latency (assuming CheckLatency is a separate function)
	if CheckLatency(pkt.AvgRtt, pkt.RTT) {
		pkt.AdditionalInfo = "High_Latency"
	}

	return pkt, nil
}

// statusAllowed returns true if code falls within any allowed [Lower,Upper] range.
// If no ranges are provided, it falls back to treating 2xx/3xx as success (to match old behavior).
func statusAllowed(code int, allowed []HttpStatusCode) bool {
	if len(allowed) == 0 {
		return code >= 200 && code < 400
	}
	for _, r := range allowed {
		low, high := r.LowerCode, r.UpperCode
		if low > high {
			low, high = high, low
		}
		if code >= low && code <= high {
			return true
		}
	}
	return false
}
