package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

func IsWebSocketUpgrade(header *http.Header) bool {
	// websocket.IsWebSocketUpgrade(req) 
	return header.Get("Connection") == "Upgrade" && header.Get("Upgrade") == "websocket"
}

func newWebsocketProxy(proxyPass string ) http.HandlerFunc {
	target, _ := url.Parse(proxyPass)
	address := target.Host
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := net.Dial("tcp", address)
		if err != nil {
			http.Error(w, "Error contacting backend server.", 500)
			log.Printf("Error dialing websocket backend %s: %v", address, err)
			return
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Not a hijacker?", 500)
			return
		}
		nc, _, err := hj.Hijack()
		if err != nil {
			log.Printf("Hijack error: %v", err)
			return
		}
		defer nc.Close()
		defer d.Close()

		err = r.Write(d)
		if err != nil {
			log.Printf("Error copying request to target: %v", err)
			return
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(d, nc)
		go cp(nc, d)
		<-errc
	})
}
