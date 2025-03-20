package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type reverseProxy struct{}

func (n reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Scheme == "https" || r.URL.Scheme == "" {
		handleTunneling(w, r)
	} else {
		result, err := proxyGetRequest(r)
		if err != nil {
			fmt.Println(err.Error())
		}
		if &result != nil {
			fmt.Println("result status code is " + strconv.Itoa(result.StatusCode))
			w.WriteHeader(result.StatusCode)
			body, err := io.ReadAll(result.Body)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, _ = w.Write(body)
			err = result.Body.Close()
			if err != nil {
				return
			}
		} else {
			fmt.Println("result is nil")
			w.WriteHeader(503)
		}
	}
}
