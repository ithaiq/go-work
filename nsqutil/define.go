package nsqutil

import (
	"github.com/nsqio/go-nsq"
)

type NsqNodeInfo struct {
	RemoteAddr    string   `json:"remote_address"`
	Hostname      string   `json:"hostname"`
	BroadcastAddr string   `json:"broadcast_address"`
	TcpPort       uint32   `json:"tcp_port"`
	HttpPort      uint32   `json:"http_port"`
	Version       string   `json:"version"`
	Tombstones    []bool   `json:"tombstones"`
	Topics        []string `json:"topics"`
}

type Nodes struct {
	Producers []NsqNodeInfo `json:"producers"`
}

type NsqdClient struct {
	cli      *nsq.Producer
	badCount uint32
	nsqdAddr string
}
