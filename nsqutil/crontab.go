package nsqutil

import (
	log "github.com/cihub/seelog"
	"time"
)

//定时刷新nsqd client配置
type OrderCrontab struct {
	interval uint32
	ProBussi *ProducerBussi
}

func NewOrderCrontab(time uint32, cli *ProducerBussi) *OrderCrontab {
	cron := &OrderCrontab{
		interval: time,
		ProBussi: cli,
	}
	return cron
}

func (cron *OrderCrontab) process() {
	for {

		time.Sleep(time.Duration(cron.interval) * time.Second)
		cron.reload()
	}
}

func (cron *OrderCrontab) Start() {
	go cron.process()
}

//该方法只处理新增的节点,失败节点的删除逻辑 由 ProducerBussi 自己处理
func (cron *OrderCrontab) reload() error {
	if cron.ProBussi == nil {
		log.Error("[OrderCrontab] probussi is nil")
		return nil
	}
	return cron.ProBussi.FlushNsqd()
}
