package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	defer killAllProxies()

	go startHttpProxy()
	go startHttpsProxy()
	time.Sleep(time.Second)

	fmt.Println("Enter to Quit: ")
	var first string
	_, _ = fmt.Scanln(&first)
}

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
