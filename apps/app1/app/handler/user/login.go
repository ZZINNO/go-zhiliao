package user

import (
	"github.com/ZZINNO/go-zhiliao-common/logs"
	"github.com/ZZINNO/go-zhiliao-common/network"
	"github.com/ZZINNO/go-zhiliao/apps/app1/app/handler/util/user"
	"github.com/gin-gonic/gin"
)

func (Self *ApiHandler) Test(c *gin.Context) {
	logs.Debug("测试debug")
	logs.Info("测试info")
	logs.Warn("测试warn")
	logs.Error("测试debug", nil)
}

func (Self *ApiHandler) CreateUser(c *gin.Context) {
	//handler定义的service输入方法
	createReq := user.CreateUserInp{}
	if err := c.ShouldBindJSON(&createReq); err != nil {
		network.GetFailResp(err, nil).Write_resp(c)
		return
	}
	//调用handler相应的函数
	out, err := user.CreateUserService(createReq)
	if err != nil {
		network.GetFailResp(err, nil).Write_resp(c)
		return
	}
	//调用输出
	out.Resp.Write_resp(c)
	return
}

func (Self *ApiHandler) LoginByUser(c *gin.Context) {
	loginReq := user.LoginUserInp{}
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		network.GetFailResp(err, nil).Write_resp(c)
		return
	}
	out, err := user.LoginUserService(loginReq)
	if err != nil {
		network.GetFailResp(err, nil).Write_resp(c)
		return
	}
	out.Resp.Write_resp(c)
	return
}
