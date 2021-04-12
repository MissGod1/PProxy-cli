package filter

import (
	"AwesomeProject/common"
	"AwesomeProject/common/log"
	"AwesomeProject/utils"
	"github.com/imgk/divert-go"
)

func AppFilter(apps []string, session chan<- common.Session)  {
	//filter := "outbound and !loopback and !ipv6 and (tcp or udp) and event == CONNECT"
	filter := "outbound and !loopback and !ipv6 and udp and event == CONNECT"
	hd, err := divert.Open(filter, divert.LayerSocket, divert.PriorityDefault, divert.FlagSniff|divert.FlagRecvOnly)
	if err != nil {
		log.Fatalf("App Filter divert.Open err : %v", err)
	}
	log.Infof("Start App Filter: %v", filter)

	buf := make([]byte, 1)
	addr := divert.Address{}

	for  {
		_, err := hd.Recv(buf, &addr)
		if err != nil{
			log.Errorf("AppFilter hd.Recv err: %v", err)
			continue
		}
		name := utils.QueryProcessName(addr.Socket().ProcessID)
		if isExist(apps, name) {
			session <- *common.NewSession(addr)
		}
	}
}


//检测是否是需要过滤的应用
func isExist(apps []string, name string) bool {
	for _, v := range apps {
		if v == name {
			return true
		}
	}
	return false
}