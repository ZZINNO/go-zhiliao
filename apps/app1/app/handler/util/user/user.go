package user

import (
	"errors"
	"github.com/ZZINNO/go-zhiliao-common/lib/jwt"
	"github.com/ZZINNO/go-zhiliao-common/network"
	"github.com/ZZINNO/go-zhiliao/apps/app1/app/model"
	"github.com/ZZINNO/go-zhiliao/apps/app1/config"
)

type ServiceHandler struct {
}

//如果需要校验的,调用validator 并且写上相关的tag
type CreateUserInp struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

//自己定义的输出
type CreateUserOut struct {
	Resp network.ApiHandlerBaseResp
}

func CreateUserService(inp CreateUserInp) (CreateUserOut, error) {
	resp := CreateUserOut{}
	ins := model.SchemaTableUser{
		Name: inp.Name,
		Pass: inp.Pass,
	}
	if id, err := config.MysqlConnection.Table("user").Insert(ins); err != nil {
		return resp, err
	} else {
		resp.Resp = network.GetSuccessResp(map[string]interface{}{"id": id})
		return resp, nil
	}
}

type LoginUserInp struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type LoginUserOut struct {
	Resp network.ApiHandlerBaseResp
}

func LoginUserService(inp LoginUserInp) (LoginUserOut, error) {
	user := model.SchemaTableUser{}
	resp := LoginUserOut{}
	exist, err := config.MysqlConnection.Table("user").Where("name=?", inp.Name).Get(&user)
	if err != nil {
		return resp, err
	}
	if !exist {
		return resp, errors.New("用户不存在")
	}
	if user.Pass != inp.Pass {
		return resp, errors.New("用户不存在")
	}
	//做一些其他的事情
	token, err := jwt.GetAuthInstance().GenerateToken(user.Name)
	if err != nil {
		return resp, err
	}
	resp.Resp = network.GetSuccessResp(map[string]interface{}{
		"access_token": token.GetAccessToken(),
		"token_type":   "JWT",
		"expires_at":   token.GetExpiresAt(),
	})
	return resp, nil
}
