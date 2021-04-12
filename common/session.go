package common

import (
	"AwesomeProject/utils"
	"fmt"
	"github.com/imgk/divert-go"
)

//Session
type Session struct {
	localAddr  string
	localPort  uint16
	remoteAddr string
	remotePort uint16
	protocol string
}

func layer(protocol uint8) string {
	switch protocol {
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	default:
		return "Unknown"
	}
}

//new Session
func NewSession(addr divert.Address) *Session {
	return &Session{
		localAddr:  utils.ParseIPv4Address(addr.Socket().LocalAddress),
		localPort:  addr.Socket().LocalPort,
		remoteAddr: utils.ParseIPv4Address(addr.Socket().RemoteAddress),
		remotePort: addr.Socket().RemotePort,
		protocol: layer(addr.Socket().Protocol),
	}
}

//new Session
func NewSession0(buf []byte) *Session {
	localAddr := fmt.Sprintf("%d.%d.%d.%d", buf[12], buf[13], buf[14], buf[15])
	remoteAddr := fmt.Sprintf("%d.%d.%d.%d", buf[16], buf[17], buf[18], buf[19])
	//ip报文头长度
	headerLen := uint(buf[0] & 0xf) * 4
	//TCP UDP报文偏移
	offset := headerLen - 20
	localPort := (uint16(buf[20 + offset]) << 8) | uint16(buf[21 + offset])
	remotePort := (uint16(buf[22 + offset]) << 8) | uint16(buf[23 + offset])
	return &Session{
		localAddr:  localAddr,
		localPort:  localPort,
		remoteAddr: remoteAddr,
		remotePort: remotePort,
		protocol: layer(buf[9]),
	}
}

func (s *Session) String() string {
	return fmt.Sprintf("%s:%d ===> %s:%d [%s]", s.localAddr, s.localPort,
		s.remoteAddr, s.remotePort, s.protocol)
}
