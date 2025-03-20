package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	ProxyPort := -1
	for ProxyPort == -1 {
		tmp := rand.Intn(*proxyCount)
		if proxies[tmp].Status == "online" {
			ProxyPort = proxies[tmp].Port
		} else {
			fmt.Println("Proxy ", tmp, " is ", proxies[tmp].Status, " (offline)")
		}
	}

	proxyAddress := "127.0.0.1:" + strconv.Itoa(ProxyPort)
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	destConn, err := dialer.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
