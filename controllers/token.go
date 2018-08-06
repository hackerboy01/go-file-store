package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"go-file-store/vo"
	"go-file-store/models"
)

type TokenController struct {
	beego.Controller
}

// Post	根据传入的client_id 和 client_secret 创建Token
func (ctx *TokenController) Post()  {
	// 获取传入创建token的参数
	clientId := strings.Trim(ctx.GetString("client_id"), " ")
	clientSecret := strings.Trim(ctx.GetString("client_secret"), " ")
	expires, err := ctx.GetInt("expires", 300)
	if err != nil || expires < 0 {
		ctx.Ctx.Output.SetStatus(400)
		ctx.StopRun()
	}

	response := &vo.ResponseMessage{
		Code: vo.TokenCreateSuccess,
		Message: "Token 创建成功",
	}
	client, err := models.ValidateClient(clientId, clientSecret)
	ctx.Data["json"] = response
	if err != nil {
		response.Code = vo.WrongClientOdAndSecret
		response.Message = "client_id 与 client_secret 不匹配"
	} else {
		token, err := models.NewTokenAndSave(uint(expires), client)
		if err == nil {
			response.Data = make(map[string]interface{})
			response.Data["token"] = token
			response.Data["expires"] = expires
			response.Data["client_id"] = token.Client.ClientId
			response.Data["created_at"] = token.CreatedAt.Format("2006-01-02 15:04:05")
		} else {
			response.Code = vo.TokenCreateFailed
			response.Message = "未知异常，请联系管理员"
		}
	}
	ctx.ServeJSON()
}
