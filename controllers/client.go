package controllers

import (
	"github.com/astaxie/beego"
	"go-file-store/models"
	"go-file-store/vo"
)

type ClientController struct {
	beego.Controller
}

// Post	创建新的账户
func (ctx *ClientController) Post()  {
	client, err := models.NewClient()
	response := &vo.ResponseMessage{}
	response.Code = vo.ClientCreateFailed
	response.Message = "创建账户错误"
	if err == nil {
		originPassword := client.ClientSecret
		_, err = client.Save()
		if err == nil {
			response.Code = vo.ClientCreateSuccess
			response.Message = "创建账户成功"
			response.Data = make(map[string]interface{})
			response.Data["client_id"] = client.ClientId
			response.Data["client_secret"] = originPassword
		}
		ctx.Data["json"] = response
	}
	ctx.ServeJSON()
}
