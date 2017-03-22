package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// Pxy is our main struct for proxy releated attributes and methods
type Pxy struct {
	// The transport used to send proxy requests to actual server.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper
}

// NewProxy returns a new Pxy object
func NewProxy() *Pxy {
	return &Pxy{}
}

// ServeHTTP is the main handler for all requests.
func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n",
		req.Method,
		req.Host,
		req.RemoteAddr,
	)

	transport := p.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// copy the origin request, and modify according to proxy
	// standard and user rules.
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	// Set `x-Forwarded-For` header.
	// `X-Forwarded-For` contains a list of servers delimited by comma and space
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// send the modified request and get response
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// write response back to client, including status code, header and body

	for key, value := range res.Header {
		// Some header item can contains many values
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func main() {
	proxy := NewProxy()

	fmt.Println("Serve on :8080")

	http.Handle("/", proxy)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
