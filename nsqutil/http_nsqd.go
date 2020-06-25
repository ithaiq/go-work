package nsqutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type HttpBussi struct {
}

func NewHttpBussi() (*HttpBussi, error) {
	return &HttpBussi{}, nil
}

func (s *HttpBussi) getNodesIps(ip string) ([]string, error) {

	nodeUrl := fmt.Sprintf("http://%s/nodes", ip)

	resp, body, errs := gorequest.New().Timeout(2000 * time.Millisecond).Get(nodeUrl).End()
	if errs != nil && len(errs) > 0 && errs[0] != nil {
		var err error
		for _, err = range errs {
			fmt.Errorf("[GetNodesIps] http get nodes ip failed, %v", err)
		}
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http get nodes ip invalid , http code is:%v", resp.StatusCode)
	}

	result := Nodes{}

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}

	var nsqdsUrl []string

	for _, v := range result.Producers {
		url := fmt.Sprintf("%v:%v", v.BroadcastAddr, v.TcpPort)
		nsqdsUrl = append(nsqdsUrl, url)
	}

	return nsqdsUrl, nil

}

func (s *HttpBussi) GetNodesIps(lookupAddrs []string) ([]string, error) {
	for _, v := range lookupAddrs {
		nodes, err := s.getNodesIps(v)
		if err == nil {
			return nodes, nil
		}
	}
	return nil, fmt.Errorf("get nsqd nodes from lookups[%+v] failed", lookupAddrs)
}
