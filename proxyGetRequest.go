package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func proxyGetRequest(req *http.Request) (*http.Response, error) {

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
		return nil, err
	}

	if req.URL.Scheme == "" {
		req.URL.Scheme = "https"
	}
	//req.Method = "GET"
	fmt.Printf("Proxy %d dialed with Proto: %s\n", ProxyPort, req.Proto)
	fmt.Printf("URLScheme: %s\nRemoteAddr: %s\nMethod: %s\n", req.URL, req.RemoteAddr, req.Method)

	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	client := &http.Client{Transport: httpTransport}

	client.Timeout = 10 * time.Second
	req.RequestURI = ""

	res, err := client.Do(req)
	fmt.Println(err)
	fmt.Println(res.StatusCode)
	return res, err

}

func proxyGetURL(url string, headers http.Header, body string, ProxyID int) []byte {
	proxyAddress := "127.0.0.1:" + strconv.Itoa(*proxyStartPort+ProxyID)
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return []byte{}
	}

	httpTransport := &http.Transport{}
	client := &http.Client{Transport: httpTransport}

	httpTransport.Dial = dialer.Dial

	request, err := http.NewRequest(
		"GET",
		url,
		strings.NewReader(body),
	)
	if err != nil {
		return []byte{}
	}

	if headers != nil {
		request.Header = headers
	} else {
		request.Header.Set("content-type", `application/json`)
		request.Header.Set("referer", url)
		request.Header.Set("upgrade-insecure-requests", `1`)
		request.Header.Set(
			"accept",
			`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		)
		request.Header.Set(
			"user-agent",
			`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4472.124 Safari/537.36`,
		)
		request.Header.Set(
			"sec-ch-ua",
			`" Not;A Brand";v="99", "Google Chrome";v="93", "Chromium";v="93"`,
		)
	}

	response, err := client.Do(request)
	if err != nil {
		return []byte{}
	}

	bodyBytes, err := io.ReadAll(response.Body)
	return bodyBytes
}
