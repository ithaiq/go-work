package rpcserver

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"yiproto/yichat/gateway_rpc"
	"yiproto/yichat/msg"

	"github.com/cihub/seelog"
	"github.com/keonjeo/protobuf/proto"

	log "yimcom/seelog"
)

// Request 实现GatewayRpcServer接口
func (s *RPCServer) Request(ctx context.Context, m *msg.Msg) (retMsg *msg.Msg, retErr error) {
	//rid := uint64(time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	//记录用户的日志
	seelog.Current.SetContext(log.CustomFormatter{RID: strconv.FormatUint(m.Uid, 10)})
	retMsg = m
	ctxYim := &Context{
		Ctx: ctx,
		M:   m,
	}
	defer func() {
		if err := recover(); err != nil {
			retMsg.Body = nil
			retMsg.Code = uint32(msg.SystemErrorCode_SYSTEM_ERROR)
			retMsg.ErrMsg = "system error"
			seelog.Errorf("panic err:%v", err)
			seelog.Errorf("==> %s\n", string(debug.Stack()))
			retErr = fmt.Errorf("%s", "系统异常,请重试")
			seelog.Errorf("retMsg:%v", retMsg)
		}
	}()
	//seelog.Infof("request, m:%v", m)
	handler, ok := mapHandler[m.Cmd]
	if !ok {
		retMsg.Code = uint32(msg.SystemErrorCode_UNKNOW_CMD_ERROR)
		retMsg.ErrMsg = "unknown cmd"
		retMsg.Body = nil
		seelog.Errorf("unkown cmd, m:%v", m)
		return
	}

	reqType := reflect.TypeOf(handler.Req).Elem()
	req := reflect.New(reqType).Interface().(proto.Message)

	rspType := reflect.TypeOf(handler.Rsp).Elem()
	rsp := reflect.New(rspType).Interface().(proto.Message)

	err := proto.Unmarshal(m.Body, req)
	if err != nil {
		retMsg.Body = nil
		retMsg.Code = uint32(msg.SystemErrorCode_PROTO_UNMARSHAL_ERROR)
		retMsg.ErrMsg = "req proto unmarshal failed"
		seelog.Errorf("%s, err:%v", m.ErrMsg, err)
		return
	}

	//seelog.Infof("*********cmd:%d, req:%v", m.Cmd, req)

	// 默认设置错误码为0
	retMsg.Cmd = handler.Cmd
	retMsg.Code = uint32(msg.SystemErrorCode_OK)
	retMsg.ErrMsg = "ok"

	// 处理函数可以修改m的Code和ErrMsg字段
	handler.Func(ctxYim, m, req, rsp)
	seelog.Infof("*********cmd:%d, rsp:%v", m.Cmd, rsp)
	rspData, err := proto.Marshal(rsp)
	if err != nil {
		retMsg.Body = nil
		retMsg.Code = uint32(msg.SystemErrorCode_PROTO_MARSHAL_ERROR)
		retMsg.ErrMsg = "proto marshal err"
		seelog.Errorf("rsp proto marshal failed, err:%v", err)
		return
	}
	retMsg.Body = rspData
	return
}

// Router 实现GatewayRpcServer接口
func (s *RPCServer) Router(stream gateway_rpc.GatewayRpc_RouterServer) error {
	return nil
}
