package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github/erastusk/canary_lambda/env"
)

var start, connect, dns, shake time.Time

// Trace http or https requsts.
func HttpTracer(l *env.EnvVariablesLoad) *httptrace.ClientTrace {
	trace := &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) {
			dns = time.Now()
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			l.Log.Printf("\n*****************************************************")
			l.Log.Printf("\nDNS Done: %v\n", time.Since(dns))
			l.Log.Printf("\n*****************************************************")
		},

		ConnectStart: func(network, addr string) {
			connect = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("\n*****************************************************\n")
			l.Log.Printf(
				"\nConnect time: %v\nNetwork : %v\nNetwork Addr : %v",
				time.Since(connect),
				network,
				addr,
			)
		},

		GotFirstResponseByte: func() {
			l.Log.Printf("\nTime from start to first byte read: %v", time.Since(start))
			fmt.Printf("*****************************************************\n")
		},
		TLSHandshakeStart: func() {
			shake = time.Now()
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			l.Log.Printf("\nTime from start to finish TLS handshake: %v", time.Since(shake))
		},
	}
	return trace
}

// Log trace results
func LogRequest(trace httptrace.ClientTrace, e *env.EnvVariablesLoad) {
	req, _ := http.NewRequest("GET", e.Url, nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &trace))
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		fmt.Printf("\n*****************************************************\n")
		e.Log.Println(err)
		fmt.Printf("*****************************************************\n")
	}
}
