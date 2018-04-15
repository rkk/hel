package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy creates a reverse proxy for the given target.
func NewProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = "/structure/transit.ashx/status"
	}
	return &httputil.ReverseProxy{Director: director}
}

func pseudoMain() {
	proxy := NewProxy(&url.URL{
		Scheme: "http",
		Host:   "m.dk",
	})

	http.ListenAndServe(":9090", proxy)
}
