package rpcserver

import (
	"fmt"
	"net"
	"runtime/debug"
	"yimcom/seelog"

	"github.com/astaxie/beego/orm"
)

// TxEnd 事务结束函数
func TxEnd(ctx *Context, o orm.Ormer, err *error) {
	if panic := recover(); panic != nil {
		seelog.Errorf("panic err:%v", panic)
		seelog.Errorf("==> %s\n", string(debug.Stack()))
		o.Rollback()
		*err = fmt.Errorf("panic")
	}
	if *err != nil {
		seelog.Errorf("TxEnd err:%v", err)
		o.Rollback()
	} else {
		o.Commit()
	}
}

// LocalAddr 本机ip
func LocalAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if IPNet, ok := address.(*net.IPNet); ok && !IPNet.IP.IsLoopback() {
			if IPNet.IP.To4() != nil {
				return IPNet.IP.String(), nil
			}
		}
	}

	return "127.0.0.1", nil
}
