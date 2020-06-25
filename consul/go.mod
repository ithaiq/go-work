module yimcom/consul

go 1.13

require (
	github.com/hashicorp/consul/api v1.1.1-0.20190815210430-23cf22960af0
	github.com/stretchr/testify v1.4.1-0.20191213072910-41d0ae8564c6
	yimcom/rpcserver v0.0.0-00010101000000-000000000000
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
