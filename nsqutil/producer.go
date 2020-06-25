package nsqutil

import (
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/nsqio/go-nsq"
)

type ProducerBussi struct {
	sCliLock     *sync.Mutex
	nsqdCliqueue *Queue
	//保存lookupaddr地址
	lookupAddrs []string
	httpCli     *HttpBussi
}

var (
	minRetryTimes int = 2
)

func NewProducerBussi(lookupUrls []string) (*ProducerBussi, error) {

	q, cli := makeClients(lookupUrls)
	if q == nil || q.Length() == 0 || cli == nil {
		return nil, fmt.Errorf("connnetd to nsqd failed, urls:[%v]", lookupUrls)
	}

	log.Infof("now using nsqd client is [%+v] \n", cli)

	return &ProducerBussi{
		sCliLock:     new(sync.Mutex),
		nsqdCliqueue: q,
		lookupAddrs:  lookupUrls,
		httpCli:      cli,
	}, nil
}

func (s *ProducerBussi) Publish(topic string, message []byte) error {
	ins := s.getNsqdClient()
	if ins == nil {
		return fmt.Errorf("get client failed, instance is empty")
	}

	var err error
	err = ins.cli.Publish(topic, message)

	//publish失败时重试
	if err != nil {
		ins.badCount++
		return s.doRetryPublish(topic, message)
	}

	return nil

}

//定时刷新nsqd
func (s *ProducerBussi) FlushNsqd() error {
	//获取可用的nsqd集群
	nsqdAddrs, err := s.httpCli.GetNodesIps(s.lookupAddrs)
	if err != nil {
		return fmt.Errorf("get nsqd address failed, lookup addrs:[%v]", nsqdAddrs)
	}
	//考虑所有的nsqd都不可用是怎么办
	if nsqdAddrs == nil || len(nsqdAddrs) == 0 {
		return fmt.Errorf("nsqd address is empty, lookup addrs:[%v]", nsqdAddrs)
	}

	config := nsq.NewConfig()
	for _, url := range nsqdAddrs {
		//已经存在可用的nsqd链接则跳过
		if s.nsqdCliqueue.Find(url) {
			continue
		}

		pCli, err := nsq.NewProducer(url, config) // 新建生产者
		if err != nil {
			log.Errorf("new client for url:[%v] failed", url)
			continue
		}

		var st NsqdClient
		st.cli = pCli
		st.badCount = 0
		st.nsqdAddr = url

		s.sCliLock.Lock()
		s.nsqdCliqueue.Add(&st)
		s.sCliLock.Unlock()

	}
	log.Infof("---[FlushNsqd]--- time:[%v] result:[%+v] \n", FormatUnixTimestamp(), s.nsqdCliqueue.PrintClients())
	return nil

}

/***以下为内部方法***/
func (s *ProducerBussi) getNsqdClient() *NsqdClient {
	s.sCliLock.Lock()
	defer s.sCliLock.Unlock()

	for {

		instance := s.nsqdCliqueue.Get()
		if instance != nil {
			//从首部删除 往后放
			s.nsqdCliqueue.Remove()
			//错误超过一定次数直接删除不添加到尾部
			if instance.badCount >= 3 {
				log.Infof("[getNsqdClient] ----- badCount:%v, addr:%v", instance.badCount, instance.nsqdAddr)
				continue
			}

			s.nsqdCliqueue.Add(instance)
			return instance
		} else {
			log.Infof("the nsq client queue is empty \n")
			s.nsqdCliqueue.Remove()
		}

		if s.nsqdCliqueue.Length() == 0 {
			return nil
		}
	}

	return nil
}

func (s *ProducerBussi) doRetryPublish(topic string, message []byte) error {
	for index := 0; index < minRetryTimes; index++ {
		st := s.getNsqdClient()
		if st == nil {
			return fmt.Errorf("get client failed, client is empty")
		}

		err := st.cli.Publish(topic, message)
		if err == nil {
			log.Infof("[doRetryPublish] retry success instance:[%+v] message:[%+v] \n", st.cli, message)
			return nil
		}

	}
	return fmt.Errorf("retry failed to publish topic:[%v], message:[%v]", topic, message)

}

/**函数*/
//从lookup获取nsqd节点地址
func makeClients(lookupUrls []string) (*Queue, *HttpBussi) {

	httpCli, _ := NewHttpBussi()
	nsqdAddrs, err := httpCli.GetNodesIps(lookupUrls)
	if err != nil {
		log.Errorf("get nsqd address failed, lookup addrs:[%v]", lookupUrls)
		return nil, nil
	}
	if nsqdAddrs == nil || len(nsqdAddrs) == 0 {
		log.Errorf("nsqd address is empty, lookup addrs:[%v]", lookupUrls)
		return nil, nil
	}

	que := NewQueue()

	config := nsq.NewConfig()

	for _, url := range nsqdAddrs {
		pCli, err := nsq.NewProducer(url, config) // 新建生产者
		if err != nil {
			log.Errorf("new client for url:[%v] failed", url)
			continue
		}
		var st NsqdClient
		st.cli = pCli
		st.badCount = 0
		st.nsqdAddr = url

		que.Add(&st)
	}
	return que, httpCli
}
