package filter

import (
	"AwesomeProject/common"
	"AwesomeProject/common/log"
	"github.com/imgk/divert-go"
)

func ForwardFilter(address divert.Address, bBuf <-chan []byte)  {
	filter := "false"
	hd, err := divert.Open(filter, divert.LayerNetwork, divert.PriorityDefault - 2, divert.FlagSendOnly)
	if err != nil {
		log.Fatalf("Forward Filter divert.Open err: %v", err)
	}
	log.Infof("Start Forward Filter: %v", filter)
	address.Flags |= uint8(3) << 6
	//buf := make([]byte, 2048)
	for {
		buf := <- bBuf
		_, err := hd.Send(buf, &address)
		if err != nil {
			log.Errorf("Forward Filter hd.Send err: %v", err)
			continue
		}
		go log.Infof("[Recv] %v", common.NewSession0(buf).String(), len(buf))
	}
}
