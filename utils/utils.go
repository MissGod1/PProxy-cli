package utils

import (
	"PProxy-cli/common/log"
	"fmt"
	"github.com/imgk/divert-go"
	"golang.org/x/sys/windows"
	"net"
	"path/filepath"
)

//解析ipv4地址
func ParseIPv4Address(addr [16]uint8) string {
	return fmt.Sprintf("%d.%d.%d.%d", addr[3], addr[2], addr[1], addr[0])
}
func GetRemoteAddress(packet []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", packet[16], packet[17], packet[18], packet[19])
}

//根据pid查询应用名
func QueryProcessName(pid uint32) string {
	if pid <= 4 {
		return "System"
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "unknown"
	}
	defer windows.CloseHandle(h)

	var name [windows.MAX_PATH]uint16
	var n = uint32(windows.MAX_PATH)
	windows.QueryFullProcessImageName(h, 0, &name[0], &n)
	_, file := filepath.Split(windows.UTF16ToString(name[:n]))
	return file
}

//获取一个addr
func GetInterfaceIndex() divert.Address {
	hd, err := divert.Open("inbound and !loopback and udp.SrcPort == 53", divert.LayerNetwork, divert.PriorityDefault, divert.FlagSniff)
	if err != nil {
		log.Fatalf("hd open err: %v", err)
	}
	defer hd.Close()

	buf := make([]byte, 1500)
	addr := divert.Address{}
	state := make(chan struct{}, 1)
	go func() {
		for {
			_, err := hd.Recv(buf, &addr)
			if err != nil {
				log.Errorf("recv err: %v", err)
				continue
			}
			// log.Println(addr.Network().InterfaceIndex, addr.Network().SubInterfaceIndex)
			state <- struct{}{}
			break
		}
	}()
	for len(state) == 0 {
		net.LookupIP("www.baidu.com")
	}
	//<-state
	return addr
}
