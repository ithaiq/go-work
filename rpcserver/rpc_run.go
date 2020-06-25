package rpcserver

import (
	"net"
	"strconv"
	"yimcom/conf"
	"yimcom/consul"
	"yimcom/seelog"
	"yiproto/yichat/gateway_rpc"
	"yiproto/yichat/msg"

	"google.golang.org/grpc"
)

//RPCServer 用于实现GatewayRpcServer接口
type RPCServer struct{}

//Run Rpc服务启动函数
func Run(ctx *Context, addr string, bigCmd ...msg.BigCmd) error {
	s := grpc.NewServer()
	gateway_rpc.RegisterGatewayRpcServer(s, &RPCServer{})
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		seelog.Errorf("failed to listen ,addr:%s, err: %v", addr, err)
		return err
	}

	if consul.DefaultClient != nil && len(bigCmd) > 0 {
		ip, err := LocalAddr()
		if err != nil {
			seelog.Errorf("get local address failed err: %v", err)
			return err
		}
		port := lis.Addr().(*net.TCPAddr).Port
		if err = consul.DefaultClient.ServiceRegister(strconv.FormatInt(int64(bigCmd[0]), 10), ip, port, conf.GetConf().SvrName); err != nil {
			seelog.Errorf("consul register service failed err: %v", err)
			return err
		}
	}

	err = s.Serve(lis)
	if err != nil {
		seelog.Errorf("listen err:%v", err)
		return err
	}
	return nil
}
