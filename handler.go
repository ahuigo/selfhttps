package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type wrapResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrapResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

type Hander struct {
	proxys map[string]*ReverseProxy
}

func getHander(proxyPass []DomainProxy) *Hander {
	h := &Hander{
		proxys: make(map[string]*ReverseProxy, len(proxyPass)),
	}
	for _, dp := range proxyPass {
		if _, ok := h.proxys[dp.Domain]; ok {
			fmt.Println("\033[31mError: duplicated args domain " + dp.Domain + "\033[0m")
			os.Exit(0)
		}
		h.proxys[dp.Domain] = &ReverseProxy{
			http: newHttpProxy(dp.ProxyPass),
			ws:   newWebsocketProxy(dp.ProxyPass),
		}
	}
	return h
}

func (h *Hander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostname := strings.Split(r.Host, ":")[0]
	if proxy, ok := h.proxys[hostname]; ok {
		isWS := IsWebSocketUpgrade(&r.Header)
		if isWS {
			// p.ServeHTTP(w, r)
			proxy.ws.ServeHTTP(w, r)
		} else {
			sw := &wrapResponseWriter{ResponseWriter: w}
			proxy.http.ServeHTTP(sw, r)
		}
	} else {
		w.Write([]byte("please config: -d " + hostname + "=proxy_pass in command line\n"))
		w.WriteHeader(404)
	}
}
