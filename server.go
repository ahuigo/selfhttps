package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/ahuigo/gofnext"
)

var loadCertificate = loadCertificateRaw

func init() {
	loadCertificate = gofnext.CacheFn2Err(loadCertificateRaw)
}

func loadCertificateRaw(certFile, keyFile string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Printf("load certificate failed, err: %v\n", err)
		return nil, err
	}
	return &cert, nil
}

func createProxyServer(cleanup func()) *http.Server {
	config := GetConfig()
	handler := getHander(config.DomainProxys)
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				certPath, certKeyPath := getCertPath(info.ServerName)
				return loadCertificate(certPath, certKeyPath)
			},
		},
	}
	go func() {
		// ts = httptest.NewUnstartedServer(http.HandlerFunc(fn))
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Println(err)
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
		cleanup()
	}()
	return server
}
