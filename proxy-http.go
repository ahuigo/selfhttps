package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	http *httputil.ReverseProxy
	ws http.HandlerFunc
}

func newHttpProxy(proxyPass string) *httputil.ReverseProxy {
		target, _ := url.Parse(proxyPass)
		proxy := httputil.NewSingleHostReverseProxy(target)
		rawDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			req.Header.Set("X-Forwarded-Host", req.Host)
			rawDirector(req)
			// req.URL.Scheme = target.Scheme
			// req.URL.Host = target.Host
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			msg := fmt.Sprintf(`{"error": "failed to access %s (%v)"}`, proxyPass, err)
			w.WriteHeader(502)
			w.Write([]byte(msg))
		}
		return proxy
}