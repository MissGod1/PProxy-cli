package filter

import (
	"AwesomeProject/common"
	"AwesomeProject/common/log"
	"fmt"
	"github.com/imgk/divert-go"
)

//Packet Filter
func PacketFilter(ifIdx uint32, pBuf chan<- []byte, isExist func(session common.Session) bool) {
	//windivert handle
	//filter := fmt.Sprintf("ifIdx == %d and outbound and !loopback and ip and (tcp or udp)", ifIdx)
	filter := fmt.Sprintf("ifIdx == %d and outbound and !loopback and ip and udp", ifIdx)
	//filter :=  fmt.Sprintf("ifIdx == %d and outbound and !loopback and ip and udp.DstPort == 53", ifIdx)
	hd, err := divert.Open(filter, divert.LayerNetwork, divert.FlagDefault-1, divert.FlagDefault)
	if err != nil {
		log.Fatalf("divert.Open err = %v", err)
	}
	log.Infof("Start Packet Filter: %v", filter)

	addr := divert.Address{}

	//recv packet
	for {
		//buffer
		buf := make([]byte, 2048)
		n, err := hd.Recv(buf, &addr)
		if err != nil {
			log.Errorf("PacketFilter Packet Filter err = %v", err)
			continue
		}

		//udp test
		//pBuf <- buf[:n]

		//log.Println(buf[:n])
		s := common.NewSession0(buf)
		if isExist(*s) {
			pBuf <- buf[:n]
		} else {
			hd.Send(buf[:n], &addr)
		}
	}
}
