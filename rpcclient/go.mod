module yimcom/rpcclient

go 1.13

require (
	github.com/keonjeo/protobuf v1.0.2
	github.com/pkg/errors v0.8.2-0.20190227000051-27936f6d90f9
	github.com/processout/grpc-go-pool v1.2.1
	google.golang.org/grpc v1.28.0
	yimcom/consul v0.0.0-00010101000000-000000000000
	yimcom/rpcserver v0.0.0-00010101000000-000000000000
	yimcom/seelog v0.0.0-00010101000000-000000000000
	yiproto/yichat/gateway_rpc v0.0.0-00010101000000-000000000000
	yiproto/yichat/msg v0.0.0-00010101000000-000000000000
)

replace (
	yimcom/conf => ./../conf
	yimcom/consul => ./../consul
	yimcom/request => ./../request
	yimcom/rpcserver => ./../rpcserver
	yimcom/seelog => ./../seelog
	yimcom/yichatmodel => ./../yichatmodel
	yiproto/yichat/config => ./../../yiproto/yichat/config
	yiproto/yichat/friend => ./../../yiproto/yichat/friend
	yiproto/yichat/gateway_rpc => ./../../yiproto/yichat/gateway_rpc
	yiproto/yichat/group => ./../../yiproto/yichat/group
	yiproto/yichat/groupchat => ./../../yiproto/yichat/groupchat
	yiproto/yichat/helper => ./../../yiproto/yichat/helper
	yiproto/yichat/info => ./../../yiproto/yichat/info
	yiproto/yichat/login => ./../../yiproto/yichat/login
	yiproto/yichat/msg => ./../../yiproto/yichat/msg
	yiproto/yichat/security => ./../../yiproto/yichat/security
)
