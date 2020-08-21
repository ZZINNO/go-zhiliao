package router

import (
	"github.com/ZZINNO/go-zhiliao-common/network/middleware"
	"github.com/ZZINNO/go-zhiliao/apps/app1/app/handler/user"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func MiddleWareCfg(g *gin.Engine) {
	//注册全局中间件
	g.Use(middleware.Recovery(raven.DefaultClient, false))
	g.NoMethod(middleware.NoMethodHandler())
	g.NoRoute(middleware.NoRouteHandler())
	return
}

func RouterCfg(g *gin.Engine) {
	//注册全局路由
	api := g.Group("/api")
	user.SetRoute(api)
}
