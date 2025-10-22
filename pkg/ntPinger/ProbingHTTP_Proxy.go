package ntPinger

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
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
	if strings.Contains(e, "forbidden") {
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
