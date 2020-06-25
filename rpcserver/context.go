package rpcserver

import (
	"context"

	"yiproto/yichat/msg"
)

// Context 上下文
type Context struct {
	Ctx context.Context
	M   *msg.Msg
}
