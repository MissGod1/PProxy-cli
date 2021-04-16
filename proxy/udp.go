package proxy

import (
	"PProxy-cli/common/log"
	"errors"
	"fmt"
	core2 "github.com/eycorsican/go-tun2socks/core"
	"net"
	"strconv"
	"sync"
	"time"

	sscore "github.com/shadowsocks/go-shadowsocks2/core"
	sssocks "github.com/shadowsocks/go-shadowsocks2/socks"

	"github.com/eycorsican/go-tun2socks/core"
)

type udpHandler struct {
	sync.Mutex

	cipher     sscore.Cipher
	remoteAddr net.Addr
	conns      map[core.UDPConn]net.PacketConn
	timeout    time.Duration
}

func NewUDPHandler(server, cipher, password string, timeout time.Duration) core2.UDPConnHandler {
	ciph, err := sscore.PickCipher(cipher, []byte{}, password)
	if err != nil {
		log.Errorf("failed to pick a cipher: %v", err)
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		log.Errorf("failed to resolve udp address: %v", err)
	}

	return &udpHandler{
		cipher:     ciph,
		remoteAddr: remoteAddr,
		conns:      make(map[core.UDPConn]net.PacketConn, 16),
		timeout:    timeout,
	}
}

func (h *udpHandler) fetchUDPInput(conn core.UDPConn, input net.PacketConn) {
	buf := core.NewBytes(core.BufSize)

	defer func() {
		h.Close(conn)
		core.FreeBytes(buf)
	}()

	for {
		input.SetDeadline(time.Now().Add(h.timeout))
		n, _, err := input.ReadFrom(buf)
		if err != nil {
			log.Errorf("read remote failed: %v", err)
			return
		}

		addr := sssocks.SplitAddr(buf[:])
		resolvedAddr, err := net.ResolveUDPAddr("udp", addr.String())
		if err != nil {
			return
		}
		_, err = conn.WriteFrom(buf[int(len(addr)):n], resolvedAddr)
		if err != nil {
			log.Errorf("write local failed: %v", err)
			return
		}
	}
}

func (h *udpHandler) Connect(conn core.UDPConn, target *net.UDPAddr) error {
	pc, err := net.ListenPacket("udp", "")
	if err != nil {
		return err
	}
	pc = h.cipher.PacketConn(pc)

	h.Lock()
	h.conns[conn] = pc
	h.Unlock()
	go h.fetchUDPInput(conn, pc)
	if target != nil {
		log.Infof("new proxy connection for target: %v", target.Network(), target.String())
	}
	return nil
}

func (h *udpHandler) ReceiveTo(conn core.UDPConn, data []byte, addr *net.UDPAddr) error {
	h.Lock()
	pc, ok1 := h.conns[conn]
	h.Unlock()

	if ok1 {
		// Replace with a domain name if target address IP is a fake IP.
		dest := net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port))

		buf := append([]byte{0, 0, 0}, sssocks.ParseAddr(dest)...)
		buf = append(buf, data[:]...)
		_, err := pc.WriteTo(buf[3:], h.remoteAddr)
		if err != nil {
			h.Close(conn)
			return errors.New(fmt.Sprintf("write remote failed: %v", err))
		}
		return nil
	} else {
		h.Close(conn)
		return errors.New(fmt.Sprintf("proxy connection %v->%v does not exists", conn.LocalAddr(), addr))
	}
}

func (h *udpHandler) Close(conn core.UDPConn) {
	conn.Close()

	h.Lock()
	defer h.Unlock()

	if pc, ok := h.conns[conn]; ok {
		pc.Close()
		delete(h.conns, conn)
	}
}
