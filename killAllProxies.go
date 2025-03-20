package main

import (
	"fmt"
	"time"
)

func killAllProxies() {
	for i, p := range proxies {
		fmt.Printf("Killing proxy %00d... ", i)
		if (i+1)%4 == 0 {
			fmt.Println()
		}
		p.Stop()
		time.Sleep(time.Millisecond * 100)
	}
}
