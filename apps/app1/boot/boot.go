package boot

import (
	"fmt"
	"github.com/ZZINNO/go-zhiliao-common/database/mysql"
	"github.com/ZZINNO/go-zhiliao-common/database/redis"
	"github.com/ZZINNO/go-zhiliao-common/lib/daemon"
	"github.com/ZZINNO/go-zhiliao-common/lib/jwt"
	"github.com/ZZINNO/go-zhiliao-common/logs"
	"github.com/ZZINNO/go-zhiliao-common/network"
	"github.com/ZZINNO/go-zhiliao-common/util"
	"github.com/ZZINNO/go-zhiliao/apps/app1/cmd"
	"github.com/ZZINNO/go-zhiliao/apps/app1/config"
	"github.com/ZZINNO/go-zhiliao/apps/app1/router"
	"github.com/gin-gonic/gin"
	"github.com/xormplus/xorm"
	"log"
	"os"
	"strings"
)

func init() {
	initCmd()
	initConfig()
	initDaemon()
	initLogger(config.GlobalConfig.App1.Logging)
	initDataBase()
	initJwtOpt()
	initHttp(config.GlobalConfig.App1.Http, router.MiddleWareCfg, router.RouterCfg)
	logs.Info("启动完成")
}

//初始化cmd
func initCmd() {
	cmd.Init()
}

//初始化配置库
func initConfig() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}
}

//初始化daemon后台
func initDaemon() {
	port := util.S(strings.Split(config.GlobalConfig.App1.Http.Addr, ":")[1]).Int()
	daemon.InitDemon(port, "/var/run/gdemo", config.GlobalConfig.App1.Http.GlaceFul)
	if daemon.DaemonCtl.SigParmSwitch(cmd.SIG, cmd.DAEMON) {
		os.Exit(0)
	}
}

//初始化log库
func initLogger(lcfg logs.LoggerBaseConfig) {
	//计算logpath
	if lcfg.LogPath == "" {
		if config.GlobalConfig.App.Env == "product" {
			lcfg.LogPath = "../../../log"
		} else {
			lcfg.LogPath = "./log"
		}
	}
	s := logs.NewLoggerConfWithDefault()
	s.Set(
		logs.SetLogPath(lcfg.LogPath),
		logs.SetLogType(lcfg.LogType),
		logs.SetLogName(lcfg.LogName),
		logs.SetEnableRecordFileInfo(true),
		logs.SetLogLevel(logs.GetLogLevelByStr(lcfg.LogLevel)),
		logs.SetEnableStdout(lcfg.EnableStdout),
		logs.SetEnableFile(lcfg.EnableFile),
	)
	logs.InitLogger(s.GeConf(), config.GlobalConfig.App.Env)
	logs.InitSentry(config.GlobalConfig.App.SentryUrl)
}

//初始化数据库类
func initDataBase() {
	config.MysqlConnection = initMysql(config.GlobalConfig.App.Mysql, "default")
	config.AuthKvConnection = initRedis(config.GlobalConfig.App.Redis, config.GlobalConfig.App1.Redis, "auth")
	config.CacheConnection = initRedis(config.GlobalConfig.App.Redis, config.GlobalConfig.App1.Redis, "cache")
}

//初始化http服务
func initHttp(hcfg network.HttpBaseConfig, middleware func(gin *gin.Engine), route func(g *gin.Engine)) {
	var app *gin.Engine
	app = gin.New()
	if logs.Logger != nil {
		gin.DefaultWriter = logs.Logger.GetWriter()
		gin.DefaultErrorWriter = logs.Logger.GetWriter()
	}
	middleware(app)
	route(app)
	go func() {
		hOpts := []network.HttpApiServerOption{
			network.SetHttpAddress(hcfg.Addr),
			network.SetHttpEnv(config.GlobalConfig.App.Env),
			network.SetHttpEngine(app),
		}
		if err := network.InitApiServerDaemon(hOpts...); err != nil {
			logs.Error("初始化api服务错误", err)
		}
	}()
	logs.Info(fmt.Sprintf("gin服务已经启动在:%s ", hcfg.Addr))
}

//初始化Redis
func initRedis(rcfg redis.RedisBaseConfig, ccfg map[string]redis.RedisCustomConfig, key string) redis.Session {
	if cfg, ok := ccfg[key]; ok {
		sess, err := redis.InitSession(redis.CombineConfig(rcfg, cfg))
		if err != nil {
			logs.Error(fmt.Sprintf("redis.%s 初始化失败", key), err)
			return redis.Session{}
		} else {
			return *sess
		}
	} else {
		logs.Error(fmt.Sprintf("没找到redis.%s 的配置文件", key), nil)
		return redis.Session{}
	}
}

//初始化mysql
func initMysql(mcfg map[string]mysql.MysqlBaseConfig, key string) xorm.Engine {
	if cfg, ok := mcfg[key]; ok {
		cfg.ShowSql = config.GlobalConfig.App1.ShowSql
		cfg.LogLevel = config.GlobalConfig.App1.SqlLogLevel
		sess, err := mysql.InitMysqlEngine(cfg)
		if err != nil {
			logs.Error(fmt.Sprintf("mysql.%s 初始化失败", key), err)
			return xorm.Engine{}
		} else {
			return *sess
		}
	} else {
		logs.Error(fmt.Sprintf("没找到mysql.%s 的配置文件", key), nil)
		return xorm.Engine{}
	}
}

//初始化JWT
func initJwtOpt() {
	jwt.SetJwtInstanceOption(config.GlobalConfig.App.Jwt, config.AuthKvConnection)
}
