package filter

import (
	"PProxy-cli/common"
	"PProxy-cli/common/log"
	"PProxy-cli/utils"

	"github.com/imgk/divert-go"
)

//Packet Filter
func PacketFilter(filter string, address divert.Address, pBuf chan<- []byte, bBuf <-chan []byte, isExist func(ip string) bool) {
	//windivert handle

	hd, err := divert.Open(filter, divert.LayerNetwork, divert.PriorityDefault, divert.FlagDefault)
	if err != nil {
		log.Fatalf("divert.Open err = %v", err)
	}
	log.Infof("Start Packet Filter: %v", filter)

	//send to local packet
	go func() {
		address.Flags |= uint8(3) << 6
		buf := make([]byte, 2048)
		zero := make([]byte, 2048)
		for {
			copy(buf, zero) // 清零
			n := copy(buf, <-bBuf)
			_, err := hd.Send(buf, &address)
			if err != nil {
				log.Errorf("PacketFilter Filter hd.Send err: %v", err)
				continue
			}
			go log.Infof("[Recv] %v", common.NewSession0(buf).String(), n)
		}
	}()

	//recv packet
	//buffer
	buf := make([]byte, 2048)
	addr := divert.Address{}
	for {
		n, err := hd.Recv(buf, &addr)
		if err != nil {
			log.Errorf("PacketFilter Packet Filter err = %v", err)
			continue
		}

		ip := utils.GetRemoteAddress(buf)
		if isExist(ip) {
			pBuf <- buf[:n]
		} else {
			hd.Send(buf[:n], &addr)
		}
	}
}
