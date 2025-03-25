package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func startHttpProxy() {
	fmt.Println("Starting HTTP Proxy")
	rp := new(reverseProxy)
	err := http.ListenAndServe(":"+strconv.Itoa(*httpPort), rp)
	if err != nil {
		panic(err.Error())
	}
}

func startHttpsProxy() {
	fmt.Println("Starting HTTPS Proxy")
	rp := new(reverseProxy)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(*httpsPort), "./certs/domain.cert.pem", "./certs/private.key.pem", rp)
	if err != nil {
		panic(err.Error())
	}
}
