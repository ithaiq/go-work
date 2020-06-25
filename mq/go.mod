module yimcom/mq

go 1.13

require (
	github.com/streadway/amqp v0.0.0-20190815230801-eade30b20f1d
	yimcom/comm v0.0.0-00010101000000-000000000000
	yimcom/seelog v0.0.0-00010101000000-000000000000
)

replace (
	yimcom/comm => ./../comm
	yimcom/conf => ./../conf
	yimcom/consul => ./../consul
	yimcom/mysql => ./../mysql
	yimcom/nsqutil => ./../nsqutil
	yimcom/request => ./../request
	yimcom/rpcserver => ./../rpcserver
	yimcom/seelog => ./../seelog
	yimcom/yichatmodel => ./../yichatmodel
	yiproto/yichat/gateway_rpc => ./../../yiproto/yichat/gateway_rpc
	yiproto/yichat/groupchat => ./../../yiproto/yichat/groupchat
	yiproto/yichat/msg => ./../../yiproto/yichat/msg
	yiproto/yichat/yichat_topic => ./../../yiproto/yichat/yichat_topic
)
