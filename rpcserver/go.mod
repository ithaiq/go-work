module yimcom/rpcserver

go 1.13

require (
	github.com/astaxie/beego v1.11.1
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/keonjeo/protobuf v1.0.2
	google.golang.org/grpc v1.28.0
	yimcom/conf v0.0.0-00010101000000-000000000000
	yimcom/consul v0.0.0-00010101000000-000000000000
	yimcom/seelog v0.0.0-00010101000000-000000000000
	yimcom/yichatmodel v0.0.0-00010101000000-000000000000
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
	yiproto/yichat/gateway_rpc => ./../../yiproto/yichat/gateway_rpc
	yiproto/yichat/msg => ./../../yiproto/yichat/msg
)
