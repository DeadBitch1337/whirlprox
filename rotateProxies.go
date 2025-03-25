package main

import (
	"fmt"
	"time"
)

func rotateProxies() {
	nextProxy := 0
	resetTime := time.Minute
	proxiesToRotate := *proxyCount / 10
	if *proxyCount > 100 {
		resetTime = time.Second * 6
		proxiesToRotate = *proxyCount / 100
	}
	for true {
		time.Sleep(resetTime)
		for i := nextProxy; i < nextProxy+proxiesToRotate; i++ {
			targetProxy := i % *proxyCount
			fmt.Printf("Rotating proxy %d\n", targetProxy)
			proxies[targetProxy].Stop()
			time.Sleep(time.Millisecond * 100)
			go proxies[targetProxy].Start()
		}
		nextProxy = nextProxy + proxiesToRotate
	}
}
