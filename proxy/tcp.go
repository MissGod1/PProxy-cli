package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	sscore "github.com/shadowsocks/go-shadowsocks2/core"
	sssocks "github.com/shadowsocks/go-shadowsocks2/socks"

	"AwesomeProject/common/log"
	"github.com/eycorsican/go-tun2socks/core"
)

type tcpHandler struct {
	cipher  sscore.Cipher
	server  string
}

func (h *tcpHandler) handleInput(conn net.Conn, input io.ReadCloser) {
	defer func() {
		conn.Close()
		input.Close()
	}()
	io.Copy(conn, input)
}

func (h *tcpHandler) handleOutput(conn net.Conn, output io.WriteCloser) {
	defer func() {
		conn.Close()
		output.Close()
	}()
	io.Copy(output, conn)
}

func NewTCPHandler(server, cipher, password string) core.TCPConnHandler {
	ciph, err := sscore.PickCipher(cipher, []byte{}, password)
	if err != nil {
		log.Errorf("failed to pick a cipher: %v", err)
	}
	return &tcpHandler{
		cipher:  ciph,
		server:  server,
	}
}

func (h *tcpHandler) Handle(conn net.Conn, target *net.TCPAddr) error {
	if target == nil {
		log.Errorf("unexpected nil target")
	}

	// Connect the relay server.
	rc, err := net.Dial("tcp", h.server)
	if err != nil {
		return errors.New(fmt.Sprintf("dial remote server failed: %v", err))
	}
	rc = h.cipher.StreamConn(rc)

	// Replace with a domain name if target address IP is a fake IP.

	dest := net.JoinHostPort(target.IP.String(), strconv.Itoa(target.Port))

	// Write target address.
	tgt := sssocks.ParseAddr(dest)
	_, err = rc.Write(tgt)
	if err != nil {
		log.Errorf("send target address failed: %v", err)
		return fmt.Errorf("send target address failed: %v", err)
	}

	go h.handleInput(conn, rc)
	go h.handleOutput(conn, rc)

	log.Infof("new proxy connection for target: %v", target.Network(), dest)
	return nil
}
