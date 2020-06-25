package rpcserver

import (
	"strings"
	"yimcom/seelog"
	"yimcom/yichatmodel"
	"yiproto/yichat/msg"

	"github.com/astaxie/beego/orm"
	"github.com/keonjeo/protobuf/proto"
)

// Handler 处理器
type Handler struct {
	Cmd  uint32
	Req  proto.Message
	Rsp  proto.Message
	Func HandlerFunc
}

// HandlerFunc 协议处理函数
type HandlerFunc func(ctx *Context, m *msg.Msg, req proto.Message, rsp proto.Message)

var mapHandler = map[uint32]Handler{}

// Register 注册命令处理器
func Register(reqCmd uint32, req proto.Message, rspCmd uint32, rsp proto.Message, handler HandlerFunc) {
	mapHandler[reqCmd] = Handler{
		Cmd:  rspCmd,
		Req:  req,
		Rsp:  rsp,
		Func: handler,
	}
}

// GetSupportCmd 获取注册的命令字
func GetSupportCmd() []uint32 {
	cmds := []uint32{}
	for cmd := range mapHandler {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func RegisterToMysql(ctx *Context, o orm.Ormer, bigCmd uint32, svrAddr string) (err error) {
	defer func() {
		if err != nil {
			seelog.Errorf("register to mysql err:%v", err)
		}
	}()

	if strings.Index(svrAddr, ":") == 0 {
		svrAddr = "127.0.0.1" + svrAddr
	}

	rs := &yichatmodel.RpcService{
		BigCmd:  bigCmd,
		SvrAddr: svrAddr,
		Status:  1,
	}

	_, err = o.InsertOrUpdate(rs)
	if err != nil {
		return err
	}

	return nil
}
