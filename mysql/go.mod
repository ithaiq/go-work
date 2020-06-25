module yimcom/mysql

go 1.13

require (
	github.com/astaxie/beego v1.11.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/robfig/cron v1.2.0
	yimcom/conf v0.0.0-00010101000000-000000000000
	yimcom/seelog v0.0.0-00010101000000-000000000000
	yimcom/yichatmodel v0.0.0-00010101000000-000000000000
)

replace (
	yimcom/conf => ./../conf
	yimcom/request => ./../request
	yimcom/seelog => ./../seelog
	yimcom/yichatmodel => ./../yichatmodel
)
