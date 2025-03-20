package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type SubProxy struct {
	ID             int
	IP             string
	Port           int
	Status         string
	BlockList      []string
	ControlChannel chan string

	cmd *exec.Cmd
}

func newSubProxy(id int) *SubProxy {
	return &SubProxy{
		ID:             id,
		IP:             "unknown",
		Port:           *proxyStartPort + id,
		Status:         "offline",
		BlockList:      make([]string, 0),
		ControlChannel: make(chan string),
	}
}

func (p *SubProxy) Start() {
	if p.cmd != nil {
		_ = p.cmd.Process.Kill()
	}
	var procIn io.Reader
	var procOut, procErr io.Writer
	p.cmd = exec.Command("tor.exe", "-f", "./torrc-tmp/torrc."+strconv.Itoa(p.ID))
	p.cmd.Stdin = procIn
	p.cmd.Stdout = procOut
	p.cmd.Stderr = procErr

	err := p.cmd.Start()
	if err != nil {
		fmt.Printf("proxy %d exited on Start with error: %s", p.ID, err.Error())
		log.Fatal(err.Error())
	}
}

func (p *SubProxy) UpdateIP() string {
	resultBytes := proxyGetURL("https://httpbin.org/ip", nil, "", p.ID)
	if resultBytes != nil && len(resultBytes) > 0 {
		p.Status = "online"
		var result struct {
			IP string `json:"origin"`
		}
		err := json.Unmarshal(resultBytes, &result)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			p.IP = result.IP
			if p.IP == "" {
				p.IP = "unknown"
			}
		}
	} else {
		p.IP = "unknown"
		if resultBytes == nil {
			p.Status = "offline"
		}
	}
	fmt.Printf("Proxy #%03d IP: %s\t| Status: %s\n", p.ID, p.IP, p.Status)
	return p.IP
}

func (p *SubProxy) Stop() {
	err := p.cmd.Process.Kill()
	if err != nil {
		fmt.Printf("proxy %d exited on Stop with error: %s", p.ID, err.Error())
	}
	p.Status = "offline"
}

func (p *SubProxy) Monitor() {
	for {
		msg := <-p.ControlChannel
		if msg == "stop" {
			p.Stop()
		} else if msg == "start" {
			p.Start()
		} else if msg == "restart" {
			p.Stop()
			time.Sleep(time.Millisecond * 100)
			p.Start()
		} else if msg == "update" {
			p.UpdateIP()
		}
	}
}
