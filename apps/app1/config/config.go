package config

import (
	"github.com/ZZINNO/go-zhiliao-common/cfg"
	"github.com/ZZINNO/go-zhiliao-common/database/redis"
	"github.com/ZZINNO/go-zhiliao-common/logs"
	"github.com/ZZINNO/go-zhiliao-common/network"
	"github.com/ZZINNO/go-zhiliao/apps/app1/cmd"
	"github.com/xormplus/xorm"
)

type CustomConfig struct {
	App  cfg.AppConfig
	App1 AppCustom
}

type AppCustom struct {
	Module      string
	ShowSql     bool
	SqlLogLevel string
	Http        network.HttpBaseConfig
	Logging     logs.LoggerBaseConfig
	Redis       map[string]redis.RedisCustomConfig
}

var GlobalConfig CustomConfig

//redis  链接
var CacheConnection redis.Session  //缓存链接
var AuthKvConnection redis.Session //认证用链接

//mysql链接
var MysqlConnection xorm.Engine

func InitConfig() error {
	cust := CustomConfig{}
	c := cfg.NewConfigParserDefault()
	c.ConfigName = cmd.ENV + ".config"
	c.InitConf()
	err := c.Reformat(&cust)
	if err != nil {
		return err
	}
	GlobalConfig = cust
	return nil
}
