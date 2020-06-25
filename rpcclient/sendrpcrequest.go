package rpcclient

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"yimcom/rpcserver"
	"yimcom/seelog"
	"yiproto/yichat/gateway_rpc"
	"yiproto/yichat/msg"

	"github.com/keonjeo/protobuf/proto"
	"google.golang.org/grpc"
)

const (
	// RPCTimeout RPC请求超时时间
	RPCTimeout = time.Second * 10
)

// SendRPCRequest 发送rpc请求
func SendRPCRequest(ctx *rpcserver.Context, bigCmd uint32, m *msg.Msg, timeout time.Duration) (*msg.Msg, error) {
	conn, err := Acquire(bigCmd)
	if err != nil {
		seelog.Errorf("SendRPCRequest => rpc dial err:%v", err)
		return nil, err
	}
	defer conn.Close() // 不会真正地关闭连接，仅是释放回连接池

	c := gateway_rpc.NewGatewayRpcClient(conn.ClientConn)
	ctxRPC, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rspMsg, err := c.Request(ctxRPC, m)
	if err != nil {
		seelog.Errorf("SendRPCRequest => rpc request err: %v", err)
		return nil, err
	}

	return rspMsg, nil
}

// SendRPCRequestMsg 发送RPC请求
func SendRPCRequestMsg(ctx *rpcserver.Context, bigCmd uint32, reqCmd uint32, req proto.Message,
	rspCmd uint32, rsp proto.Message) error {
	if ctx.M == nil {
		seelog.Infof("SendRPCRequestMsg => ctx.M is nil")
		ctx.M = &msg.Msg{}
	}
	newMsg := proto.Clone(ctx.M)

	rpcMsg, ok := newMsg.(*msg.Msg)
	if !ok {
		seelog.Errorf("SendRPCRequestMsg => ctx.M is nil")
	}
	rpcMsg.Cmd = reqCmd
	rpcMsg.BigCmd = bigCmd

	body, err := proto.Marshal(req)
	if err != nil {
		seelog.Errorf("SendRPCRequestMsg => req marshal err: %v", err)
		return err
	}
	rpcMsg.Body = body
	rpcRspMsg, err := SendRPCRequest(ctx, bigCmd, rpcMsg, RPCTimeout)
	if err != nil {
		seelog.Errorf("SendRPCRequestMsg => Fail to SendRPCRequest err: %v", err)
		return err
	}
	if rpcRspMsg.Code != 0 {
		err = errors.New(fmt.Sprintf("code:%d err msg:%s", rpcRspMsg.Code, rpcRspMsg.ErrMsg))
		return err
	}

	if rpcRspMsg.Code == 0 && rpcRspMsg.Cmd == rspCmd {
		err = proto.Unmarshal(rpcRspMsg.Body, rsp)
		if err != nil {
			seelog.Errorf("SendRPCRequestMsg => Fail to proto.Unmarshal, err: %v, body: %v", err, rpcRspMsg.Body)
			return err
		}
	}
	return nil
}

func SendYimgateRequestMsg(ctx *rpcserver.Context, bigCmd uint32) (*msg.Msg, error) {
	if ctx.M == nil {
		seelog.Errorf("SendYimgateRequestMsg => ctx.M is nil")
		return nil, errors.New("msg info is null")
	}
	rpcRspMsg, err := SendRPCRequest(ctx, bigCmd, ctx.M, RPCTimeout)
	if err != nil {
		seelog.Errorf("SendYimgateRequestMsg => Fail to SendRPCRequest, cmd: %v, err: %v", bigCmd, err)
		return nil, errors.Wrap(err, "SendRPCRequest error")
	}
	return rpcRspMsg, nil

}

//指定目标的短连接rpc调用
func doRPCRequest(ctx *rpcserver.Context, addr string, m *msg.Msg, timeout time.Duration) (*msg.Msg, error) {
	var err error

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		seelog.Errorf("[DoRPCRequest] rpc dial err:%v", err)
		return nil, err
	}
	defer conn.Close()
	c := gateway_rpc.NewGatewayRpcClient(conn)
	ctxRPC, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rspMsg, err := c.Request(ctxRPC, m)
	if err != nil {
		seelog.Errorf("[DoRPCRequest] rpc request err:%v", err)
		return nil, err
	}

	return rspMsg, nil
}

func DoRPCRequestMsg(ctx *rpcserver.Context, addr string, bigCmd uint32, reqCmd uint32, req proto.Message,
	rspCmd uint32, rsp proto.Message) error {
	if ctx.M == nil {
		seelog.Errorf("ctx.M is nil")
		ctx.M = &msg.Msg{}
	}
	newMsg := proto.Clone(ctx.M)
	rpcMsg, ok := newMsg.(*msg.Msg)
	if !ok {
		seelog.Errorf("ctx.M is nil")
	}
	rpcMsg.Cmd = reqCmd
	rpcMsg.BigCmd = bigCmd
	body, err := proto.Marshal(req)
	if err != nil {
		seelog.Errorf("req marshal err:%v", err)
		return err
	}
	rpcMsg.Body = body

	rpcRspMsg, err := doRPCRequest(ctx, addr, rpcMsg, RPCTimeout)
	if err != nil {
		seelog.Errorf("SendRPCRequest err:%v", err)
		return err
	}

	if rpcRspMsg.Code != 0 {
		err = errors.New(fmt.Sprintf("code:%d err msg:%s", rpcRspMsg.Code, rpcRspMsg.ErrMsg))
		return err
	}

	if rpcRspMsg.Code == 0 && rpcRspMsg.Cmd == rspCmd {
		err = proto.Unmarshal(rpcRspMsg.Body, rsp)
		if err != nil {
			seelog.Errorf("pb unmarshal err:%v, body:%v", err, rpcRspMsg.Body)
			return err
		}
	}
	return nil
}
