package main

import (
	"AwesomeProject/common"
	"AwesomeProject/common/log"
	"AwesomeProject/filter"
	"AwesomeProject/proxy"
	"AwesomeProject/utils"
	"flag"
	"fmt"
	"github.com/eycorsican/go-tun2socks/core"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	path := flag.String("config", "config.json", "Configure File")
	level := flag.String("log", "info", "Log Level")
	udp := flag.Bool("udp", false, "enable udp mode")
	tcp := flag.Bool("tcp", false, "enable tcp mode")
	timeOut := flag.Duration("timeout", 1*time.Minute, "udp timeout")
	flag.Parse()

	//set log level
	switch strings.ToLower(*level) {
	case "debug":
		log.SetLevel(log.DEBUG)
	case "info":
		log.SetLevel(log.INFO)
	case "error":
		log.SetLevel(log.ERROR)
	case "warning":
		log.SetLevel(log.WARNING)
	case "none":
		log.SetLevel(log.NONE)
	default:
		log.SetLevel(log.INFO)
	}
	config := common.NewConfig(*path)

	//register handle
	core.RegisterTCPConnHandler(proxy.NewTCPHandler(core.ParseTCPAddr(config.Server, config.Port).String(),
		config.Method, config.Password))
	core.RegisterUDPConnHandler(proxy.NewUDPHandler(core.ParseUDPAddr(config.Server, config.Port).String(),
		config.Method, config.Password, *timeOut))

	//variable define
	rip := make(chan string, 8)
	remoteMap := make(map[string]struct{}) //对端ip表
	pBuf := make(chan []byte, 8)
	bBuf := make(chan []byte, 8)
	mutex := sync.Mutex{}
	addr := utils.GetInterfaceIndex()

	stack := core.NewLWIPStack()
	core.RegisterOutputFn(func(data []byte) (int, error) {
		log.Debugf("OutputFn: %v", len(data))
		bBuf <- data
		return len(data), nil
	})
	//中转协程，处理session存储和包转发
	go func() {
		for {
			select {
			case ip := <-rip:
				mutex.Lock()
				remoteMap[ip] = struct{}{}
				mutex.Unlock()
				go log.Infof("[CONNECT] %v", ip)
			case buf := <-pBuf:
				_, err := stack.Write(buf)
				if err != nil {
					go log.Errorf("stack.Write err: %v", err)
				}
				go log.Infof("[Send] %v", common.NewSession0(buf), len(buf))
			}

		}
	}()

	//mode
	appFilter := "outbound and !loopback and !ipv6 and event == CONNECT and "
	packetFilter := fmt.Sprintf("ifIdx == %d and outbound and !loopback and ip and ", addr.Network().InterfaceIndex)
	if (*udp && *tcp) || (!*udp && !*tcp) {
		appFilter += "(tcp or udp)"
		packetFilter += "(tcp or udp)"
	} else if *udp {
		appFilter += "udp"
		packetFilter += "udp"
	} else {
		appFilter += "tcp"
		packetFilter += "tcp"
	}

	//Start Filter
	go filter.AppFilter(appFilter, config.Apps, rip)
	go filter.PacketFilter(packetFilter, addr, pBuf, bBuf, func(ip string) bool {
		mutex.Lock()
		defer mutex.Unlock()
		if _, ok := remoteMap[ip]; ok {
			return true
		} else {
			return false
		}
	})

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
