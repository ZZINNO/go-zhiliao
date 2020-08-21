package user

import (
	"github.com/ZZINNO/go-zhiliao/apps/app1/app"
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	app.ApiHandler
}

func SetRoute(r *gin.RouterGroup) {
	userHandler := ApiHandler{}
	demo := r.Group("user")
	demo.GET("/test", userHandler.Test)
	demo.POST("/add", userHandler.CreateUser)
	demo.POST("/auth/login", userHandler.LoginByUser)
}
