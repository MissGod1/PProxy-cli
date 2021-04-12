package main

import (
	"AwesomeProject/common"
	"AwesomeProject/common/log"
	"AwesomeProject/filter"
	"AwesomeProject/proxy"
	"AwesomeProject/utils"
	"flag"
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
	level := flag.String("level", "info", "Log Level")
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
		config.Method, config.Password, 1*time.Minute))

	//variable define
	session := make(chan common.Session, 8)
	sessionMap := make(map[common.Session]struct{})
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
			case s := <-session:
				mutex.Lock()
				sessionMap[s] = struct{}{}
				mutex.Unlock()
				go log.Infof("[CONNECT] %v", s.String())
			case buf := <-pBuf:
				stack.Write(buf)
				go log.Infof("[Send] %v", common.NewSession0(buf), len(buf))
			}

		}
	}()

	//Start Filter
	go filter.ForwardFilter(addr, bBuf)
	go filter.PacketFilter(addr.Network().InterfaceIndex, pBuf, func(session common.Session) bool {
		mutex.Lock()
		defer mutex.Unlock()
		if _, ok := sessionMap[session]; ok {
			return true
		} else {
			return false
		}
	})
	go filter.AppFilter(config.Apps, session)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
