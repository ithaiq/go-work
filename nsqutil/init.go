package nsqutil

import "errors"

var nsqdBussi *ProducerBussi

func Init(lookupAddrs []string) error {
	var err error
	nsqdBussi, err = NewProducerBussi(lookupAddrs)
	if err != nil {
		return err
	}

	//启动定时任务,定时检测是否有新的可用nsqd节点，并建立连接
	cronTask := NewOrderCrontab(10, nsqdBussi)
	cronTask.Start()
	return nil
}

func Publish(topic string, body []byte) error {
	if nsqdBussi == nil {
		return errors.New("you need to init nsqutil first")
	}
	return nsqdBussi.Publish(topic, body)
}
