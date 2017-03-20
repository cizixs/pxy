package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Pxy struct {
	proxy *httputil.ReverseProxy
}

func NewProxy() *Pxy {
	director := func(req *http.Request) {
		source := req.URL

		req.URL.Scheme = source.Scheme
		req.URL.Host = source.Host
	}

	return &Pxy{
		proxy: &httputil.ReverseProxy{Director: director},
	}
}

func (p *Pxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Host)
	p.proxy.ServeHTTP(w, r)
}

func main() {
	proxy := NewProxy()

	fmt.Println("Serve on :8080")

	http.Handle("/", proxy)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
