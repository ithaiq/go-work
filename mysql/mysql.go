package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"math"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"yimcom/conf"
	"yimcom/yichatmodel"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"yimcom/seelog"
)

var (
	mysqlMap     sync.Map
	mu           sync.Mutex
	appSourceAll []uint32
)

const (
	MaxAppSourceId = 2000 //hub_common数据库
	CronTime       = 5    //5分钟执行
)

func init() {
	//StartTimer 启动定时器,需要单独启动线程执行
	c := cron.New()
	f := func() {
		defaultOrm := getConnByID(MaxAppSourceId)
		if defaultOrm != nil {
			var delAppSourceList []int
			_, err := defaultOrm.Raw("select app_source from app_cfg where app_cfg.disabled = 1 order by id").QueryRows(&delAppSourceList)
			if err != nil {
				return
			}
			for _, appSource := range delAppSourceList {
				app := fmt.Sprintf("%v", appSource)
				if _, ok := mysqlMap.Load(app); ok {
					mysqlMap.Delete(app)
				}
			}
		}
	}
	timeFmt := fmt.Sprintf("0 0/%d * * * ?", CronTime)
	seelog.Info("StartTimer timer format: ", timeFmt)
	c.AddFunc(timeFmt, f)
	// 开始
	c.Start()
}

//InitMysqlDb 初始化mysql数据库
func InitMysqlDb(username string, password string, ip string, port uint32, database string,
	idleConn int, maxConn int, maxLifttime int, appSourceID interface{}) (error, orm.Ormer) {
	mu.Lock()
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("InitMysqlDb painc error:", err)
			seelog.Error(string(debug.Stack()))
		}
		mu.Unlock()
	}()

	var dbName string
	if appSourceID == MaxAppSourceId {
		dbName = "default"
	} else {
		dbName = fmt.Sprintf("%v", appSourceID)
	}

	if conn := getConnByID(appSourceID); conn != nil {
		return nil, conn
	}

	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&loc=Local", username, password, "tcp",
		ip, port, database)
	seelog.Infof("mysql dsn:%s", dsn)
	seelog.Infof("mysql idle conn:%d, max conn:%d, max lift time:%d", idleConn, maxConn, maxLifttime)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase(dbName, "mysql", dsn, idleConn, maxConn)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			o := orm.NewOrm()
			o.Using(dbName)
			mysqlMap.Store(fmt.Sprintf("%v", appSourceID), dbName)
			return nil, o
		} else {
			seelog.Errorf("get mysql db, err:%v", err)
			return err, nil
		}
	}

	orm.Debug = true
	mdb, err := orm.GetDB(dbName)
	mdb.SetConnMaxLifetime(time.Duration(maxLifttime) * time.Second)
	mdb.SetMaxIdleConns(idleConn)
	mdb.SetMaxOpenConns(maxConn)
	o := orm.NewOrm()
	o.Using(dbName)
	//mysqlMap[appSourceID] = o
	mysqlMap.Store(fmt.Sprintf("%v", appSourceID), dbName)
	return err, o
}

//获取map中mysql链接
func getConnByID(appSource interface{}) orm.Ormer {
	if v, ok := mysqlMap.Load(fmt.Sprintf("%v", appSource)); ok {
		if dbName, ok := v.(string); ok {
			o := orm.NewOrm()
			o.Using(dbName)
			return o
		}
	}
	return nil
}

//获取一个mysql链接
func GetMysqlConn(appSource interface{}) orm.Ormer {
	if appSource == nil || fmt.Sprintf("%v", appSource) == "" {
		seelog.Errorf("mysql [GetMysqlConn] appSource is %v", appSource)
		return nil
	}
	if fmt.Sprintf("%v", appSource) == "0" {
		panic("GetMysqlConn has err")
	}
	if conn := getConnByID(appSource); conn != nil {
		return conn
	}
	err, ormRet := getConnFromHubCfg(appSource)
	if err != nil {
		seelog.Errorf("[GetMysqlConn] err:%v", err)
	}
	return ormRet
}

//读取hub_common配置并初始化链接
func getConnFromHubCfg(appSource interface{}) (error, orm.Ormer) {
	var (
		cfg        yichatmodel.MysqlCfg
		defaultOrm = getConnByID(MaxAppSourceId)
		err        error
	)
	if defaultOrm == nil {
		conf := conf.GetConf()
		// 初始化 mysql
		err, defaultOrm = InitMysqlDb(conf.MysqlUser, conf.MysqlPwd, conf.MysqlHost, conf.MysqlPort,
			conf.MysqlDb, conf.MysqlMaxIdle, conf.MysqlMaxConn, conf.MysqlMaxLifetime, MaxAppSourceId)
		if err != nil {
			seelog.Errorf("[getConnFromHubCfg] has err:%v", err)
			return err, nil
		}
	}
	sql := fmt.Sprintf("SELECT `host`,`port`,`user`,pwd,db_name,max_conn,max_idle,max_life_time from app_cfg left JOIN mysql_cfg on mysql_cfg.id = app_cfg.mysql_cfg_id where app_cfg.app_source = %v and app_cfg.disabled = 0", appSource)
	err = defaultOrm.Raw(sql).QueryRow(&cfg.Host, &cfg.Port, &cfg.User, &cfg.Pwd,
		&cfg.DbName, &cfg.MaxConn, &cfg.MaxIdle, &cfg.MaxLifeTime)
	if err != nil {
		return err, nil
	}
	err, ormRet := InitMysqlDb(cfg.User, cfg.Pwd, cfg.Host, cfg.Port, cfg.DbName,
		cfg.MaxIdle, cfg.MaxConn, cfg.MaxLifeTime, appSource)
	if err != nil {
		seelog.Errorf("[getConnFromHubCfg]  err:%v", err)
		return err, ormRet
	}
	return nil, ormRet
}

func GetRedisCfgByDb(appSource interface{}) (error, *yichatmodel.RedisCfg) {
	var (
		cfg        yichatmodel.RedisCfg
		defaultOrm orm.Ormer
		err        error
	)
	defaultOrm = getConnByID(MaxAppSourceId)
	if defaultOrm == nil {
		conf := conf.GetConf()
		// 初始化 mysql
		err, defaultOrm = InitMysqlDb(conf.MysqlUser, conf.MysqlPwd, conf.MysqlHost, conf.MysqlPort,
			conf.MysqlDb, conf.MysqlMaxIdle, conf.MysqlMaxConn, conf.MysqlMaxLifetime, MaxAppSourceId)
		if err != nil {
			seelog.Errorf("[InitMysqlDb] has err:%v", err)
			return err, nil
		}
	}
	sql := fmt.Sprintf("SELECT `host`,`port`,auth,max_idle,max_active,pool_size,`db` from app_cfg left JOIN redis_cfg on redis_cfg.id = app_cfg.redis_cfg_id where app_cfg.app_source = %v and app_cfg.disabled = 0", appSource)
	err = defaultOrm.Raw(sql).QueryRow(&cfg.Host, &cfg.Port, &cfg.Auth, &cfg.MaxIdle, &cfg.MaxActive, &cfg.PoolSize, &cfg.DB)
	return err, &cfg
}

//TODO 获取所有app,传true代表获取所有的 注意mysql链接获取为空
func GetAllAppSource(isGetAll bool) (error, []uint32) {
	if len(appSourceAll) != 0 {
		return nil, appSourceAll
	}
	var (
		appSourceList []uint32
		err           error
		defaultOrm    orm.Ormer
	)
	defaultOrm = getConnByID(MaxAppSourceId)
	if defaultOrm == nil {
		conf := conf.GetConf()
		// 初始化 mysql
		err, defaultOrm = InitMysqlDb(conf.MysqlUser, conf.MysqlPwd, conf.MysqlHost, conf.MysqlPort,
			conf.MysqlDb, conf.MysqlMaxIdle, conf.MysqlMaxConn, conf.MysqlMaxLifetime, MaxAppSourceId)
		if err != nil {
			seelog.Errorf("[GetAllAppSource] InitMysqlDb has err:%v", err)
			return err, nil
		}
	}
	sql := "select app_source from app_cfg where app_cfg.disabled = 0 order by id"
	if isGetAll {
		sql = "select app_source from app_cfg order by id"
	}
	_, err = defaultOrm.Raw(sql).QueryRows(&appSourceList)
	return err, appSourceList
}

func GetOneAppSource() (error, uint32) {
	err, appSrcList := GetAllAppSource(false)
	if err != nil || len(appSrcList) == 0 {
		return errors.New("GetOneAppSource has err"), math.MaxInt16
	}
	return nil, appSrcList[0]
}

//TODO 获取所有app,传true代表获取所有的 注意mysql链接获取为空
func GetAllAppData(isGetAll bool) (error, []*yichatmodel.AppCfg) {
	var (
		data []*yichatmodel.AppCfg
		err  error
	)
	defaultOrm := getConnByID(MaxAppSourceId)
	if defaultOrm == nil {
		conf := conf.GetConf()
		// 初始化 mysql
		err, defaultOrm = InitMysqlDb(conf.MysqlUser, conf.MysqlPwd, conf.MysqlHost, conf.MysqlPort,
			conf.MysqlDb, conf.MysqlMaxIdle, conf.MysqlMaxConn, conf.MysqlMaxLifetime, MaxAppSourceId)
		if err != nil {
			seelog.Errorf("[GetAllAppData] has err:%v", err)
			return err, nil
		}
	}

	sql := "select * from app_cfg where app_cfg.disabled = 0 order by id"
	if isGetAll {
		sql = "select * from app_cfg order by id"
	}
	_, err = defaultOrm.Raw(sql).QueryRows(&data)
	return err, data
}

func GetAppSrcByUid(userID interface{}) (appSrc uint32) {
	err, appList := GetAllAppSource(false)
	if err != nil {
		return math.MaxInt16
	}
	for _, v := range appList {
		if err := GetMysqlConn(v).Raw("select app_source from user where user_id = ? and app_source = ?", userID, v).QueryRow(&appSrc); err != nil {
			continue
		} else {
			return v
		}
	}
	return math.MaxInt16
}

func GetAppSrcByGroup(groupID interface{}) (appSrc uint32) {
	err, appList := GetAllAppSource(false)
	if err != nil {
		return math.MaxInt16
	}
	for _, v := range appList {
		if err := GetMysqlConn(v).Raw("SELECT `user`.app_source FROM `group` left JOIN `user` on `group`.group_owner = `user`.user_id where `group`.group_id = ?", groupID).QueryRow(&appSrc); err != nil {
			continue
		} else {
			return v
		}
	}
	return math.MaxInt16
}

func GetAppSrcByGroupOwner(ownerID interface{}) (appSrc uint32, err error) {
	var data yichatmodel.GroupOwnerSource
	defaultOrm := getConnByID(MaxAppSourceId)
	if defaultOrm == nil {
		conf := conf.GetConf()
		// 初始化 mysql
		err, defaultOrm = InitMysqlDb(conf.MysqlUser, conf.MysqlPwd, conf.MysqlHost, conf.MysqlPort,
			conf.MysqlDb, conf.MysqlMaxIdle, conf.MysqlMaxConn, conf.MysqlMaxLifetime, MaxAppSourceId)
		if err != nil {
			seelog.Errorf("[GetAppSrcByGroupOwner] has err:%v", err)
			return 0, err
		}
	}
	err = defaultOrm.QueryTable(&yichatmodel.GroupOwnerSource{}).Filter("group_owner_id", ownerID).One(&data)
	if err != nil {
		data.AppSourceId, err = saveGroupOwnerId(ownerID)
	}
	return data.AppSourceId, err
}

func saveGroupOwnerId(ownerID interface{}) (uint32, error) {
	appSrc := GetAppSrcByUid(ownerID)
	if appSrc == math.MaxInt16 {
		return appSrc, errors.New(fmt.Sprintf("saveGroupOwnerId has error:%v", ownerID))
	}
	defaultOrm := getConnByID(MaxAppSourceId)
	var err error
	if defaultOrm == nil {
		conf := conf.GetConf()
		// 初始化 mysql
		err, defaultOrm = InitMysqlDb(conf.MysqlUser, conf.MysqlPwd, conf.MysqlHost, conf.MysqlPort,
			conf.MysqlDb, conf.MysqlMaxIdle, conf.MysqlMaxConn, conf.MysqlMaxLifetime, MaxAppSourceId)
		if err != nil {
			seelog.Errorf("[GetAppSrcByGroupOwner] has err:%v", err)
			return 0, err
		}
	}
	_, err = defaultOrm.Raw("INSERT INTO `group_owner_source` (`group_owner_id`, `app_source_id`) VALUES (?, ?)", ownerID, appSrc).Exec()
	return appSrc, err
}

func CheckMySQLAvailable(dbHost string, dbPort uint32, dbUser, dbPassword string) bool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbUser, dbPassword, dbHost, dbPort)

	fmt.Println("dsn: ", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		seelog.Errorf("CheckMySQLAvailable => Fail to Open Database, err: %v", err)
	}

	_, err = db.Exec("select version();")
	if err != nil {
		seelog.Errorf("CheckMySQLAvailable =>  Connection issue: err: %v", err)
		return false
	} else {
		return true
	}
}
