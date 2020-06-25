package rpcserver

import (
	"net/http"
	"yimcom/seelog"
)

//SetPprof 设置pprof
func SetPprof(ctx *Context, pprofAddr string) {
	seelog.Infof("pprofAddr:%s", pprofAddr)
	go func() {
		err := http.ListenAndServe(pprofAddr, nil)
		seelog.Errorf("pprof err:%v", err)
		if err != nil {
			seelog.Errorf("pprof err:%v", err)
		} else {
			seelog.Infof("pprof at:%s", pprofAddr)
		}

	}()

}
