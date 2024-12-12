package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var transport = &http.Transport{
	MaxIdleConns:      10,
	IdleConnTimeout:   30 * time.Second,
	DisableKeepAlives: false,
}

// StartReverseProxy starts a reverse proxy server
func StartReverseProxy(listenAddr string, targetAddr string) {
	target, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatalf("Failed to parse target address: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request for: %s", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Starting reverse proxy on %s, forwarding to %s", listenAddr, targetAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// StartReverseProxyTLS starts a reverse proxy server with TLS
func StartReverseProxyTLS(listenAddr, targetAddr, certFile, keyFile string) {
	target, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatalf("Failed to parse target address: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request for: %s", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Starting TLS reverse proxy on %s, forwarding to %s", listenAddr, targetAddr)
	if err := http.ListenAndServeTLS(listenAddr, certFile, keyFile, mux); err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}
