package main

import (
	"time"
)

func main() {
	defer killAllProxies()

	go startHttpProxy()
	go startHttpsProxy()
	time.Sleep(time.Second)

	rotateProxies()
}
