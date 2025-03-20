package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var proxyStartPort *int
var cmdStartPort *int
var proxyCount *int
var slowStart *bool
var httpPort *int
var httpsPort *int
var torPath *string

var proxies []*SubProxy

func init() {
	// get cmd arguments
	proxyStartPort = flag.Int("pPort", 9000, "first port to start binding proxies on")
	cmdStartPort = flag.Int("cPort", 8000, "first port to start binding control on")
	proxyCount = flag.Int("count", 10, "number of proxies to start")
	slowStart = flag.Bool("slow", false, "start proxies slowly")
	httpPort = flag.Int("http", 7000, "port to bind HTTP to")
	httpsPort = flag.Int("https", 7001, "port to bind HTTPs to")
	torPath = flag.String("https", "tor", "port to bind HTTPs to")
	flag.Parse()

	// clear old tmp files
	fmt.Println("Creating temp directory...")
	err := os.RemoveAll("./torrc-tmp")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println("Could empty temp directory: ")
		panic(err.Error())
	}
	err = os.Mkdir("./torrc-tmp", 0777)
	if err != nil {
		fmt.Println("Could not create temp directory: ")
		panic(err.Error())
	}

	// read & prepare config
	fmt.Println("Preparing config files...")
	configBytes, err := os.ReadFile("./torrc.base")
	if err != nil {
		fmt.Println("Could not read torrc.base: ")
		panic(err.Error())
	}
	socksRegex, err := regexp.Compile("\\n*\\s*SocksPort\\s*\\d+\\s*\\n*")
	controlRegex, err := regexp.Compile("\\n*\\s*ControlPort\\s*\\d+\\s*\\n*")
	dataRegex, err := regexp.Compile("\\n*\\s*DataDirectory\\s*\\S+\\s*\\n*")

	if socksRegex.Find(configBytes) == nil {
		configBytes = []byte(string(configBytes) + "\nSocksPort 0\n")
	}
	if controlRegex.Find(configBytes) == nil {
		configBytes = []byte(string(configBytes) + "\nControlPort 0\n")
	}
	if dataRegex.Find(configBytes) == nil {
		configBytes = []byte(string(configBytes) + "\nDataDirectory 0\n")
	}

	// write config files for each proxy
	for i := 0; i < *proxyCount; i++ {
		{
			proxyPort := *proxyStartPort + i
			cmdPort := *cmdStartPort + i

			tmp := socksRegex.ReplaceAll(configBytes, []byte("\nSocksPort "+strconv.Itoa(proxyPort)+"\n"))
			tmp = controlRegex.ReplaceAll(tmp, []byte("\nControlPort "+strconv.Itoa(cmdPort)+"\n"))
			config := dataRegex.ReplaceAll(tmp, []byte("\nDataDirectory ./torrc-tmp/tor-data."+strconv.Itoa(i)+"\n"))
			err := os.WriteFile("./torrc-tmp/torrc."+strconv.Itoa(i), config, 0777)
			_ = os.Mkdir("./torrc-tmp/tor-data."+strconv.Itoa(i), 0777)
			if err != nil {
				fmt.Println("Could not write torrc files: ")
				panic(err.Error())
			}
		}

	}

	fmt.Println("Starting proxies...")
	for i := 0; i < *proxyCount; i++ {
		fmt.Printf("Starting proxy %03d... ", i)
		if (i+1)%4 == 0 {
			fmt.Println()
		}
		proxies = append(proxies, newSubProxy(i))
		proxies[i].Start()
		if *slowStart {
			time.Sleep(time.Millisecond * 500)
		} else {
			time.Sleep(time.Millisecond * 250)
		}
		go proxies[i].Monitor()
	}
	fmt.Println()

	fmt.Println("Waiting for proxies to start...")
	for i := 9; i > 0; i-- {
		fmt.Printf("%02d... ", i)
		if i%5 == 0 {
			fmt.Println()
		}
		time.Sleep(time.Second)
	}
	fmt.Println()
	go heartbeat()
}

func heartbeat() {
	for i := 0; i < *proxyCount; i++ {
		proxies[i].ControlChannel <- "update"
	}
	time.Sleep(time.Second * 3)
	heartbeat()
}
