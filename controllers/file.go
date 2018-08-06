package controllers

import (
	"github.com/astaxie/beego"
	"go-file-store/vo"
	"strings"
	"go-file-store/models"
	"fmt"
	"os"
	syspath "path"
	"go-file-store/utils"
	"log"
)

type FileController struct {
	beego.Controller
}

func (ctx *FileController) Upload()  {
	response := &vo.ResponseMessage{
		Code: vo.UploadFileSuccess,
		Message: "文件上传成功",
	}

	ctx.Data["json"] = response

	// 只允许POST方法提交
	if ctx.Ctx.Input.Method() != "POST" {
		ctx.Ctx.Output.SetStatus(405)
		response.Code = vo.RequestError
		response.Message = "请使用POST方法上传"
		ctx.ServeJSON()
		ctx.StopRun()
	}

	// 没有上传的文件
	if ctx.Ctx.Input.IsUpload() == false {
		ctx.Ctx.Output.SetStatus(400)
		response.Code = vo.RequestError
		response.Message = "文件未找到"
		ctx.ServeJSON()
		ctx.StopRun()
	}

	path := strings.Trim(ctx.GetString("path", ""), "/")
	token := strings.Trim(ctx.GetString("token"), " ")
	tokenModel, err := models.IsTokenValidate(token)

	if err == nil {
		response.Code = vo.InvalidToken
		response.Message = "非法的Token"
		ctx.ServeJSON()
		ctx.StopRun()
	}

	if tokenModel.IsTokenExpire() {
		response.Code = vo.ExpiredToken
		response.Message = "Token已过期"
		ctx.ServeJSON()
		ctx.StopRun()
	}

	file, fileHeader, err := ctx.GetFile("file")
	if err != nil {
		response.Code = vo.UploadFileFailed
		response.Message = "读取文件错误"
		ctx.ServeJSON()
		ctx.StopRun()
	}
	completeFileSaveDir := fmt.Sprintf("%s/%s", tokenModel.Client.RootDir(), path)
	completeFileSaveDir = strings.TrimRight(completeFileSaveDir, "/")
	if _, err = os.Stat(completeFileSaveDir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(completeFileSaveDir, os.ModePerm)
	}
	uploadFilePath := fmt.Sprintf("%s/%s", path, fileHeader.Filename)
	completeFilePath := fmt.Sprintf("%s/%s", completeFileSaveDir, fileHeader.Filename)
	localFilePath := fmt.Sprintf("%s/%ws", path, fileHeader.Filename)
	uploadFileExt := syspath.Ext(fileHeader.Filename)

	for {
		if _, err := os.Stat(completeFilePath); err != nil && os.IsNotExist(err) {
			break
		}
		newFileName := fmt.Sprintf("%s%s", utils.MD5(utils.GenerateRandomString(32)), uploadFileExt)
		completeFilePath = fmt.Sprintf("%s/%s", completeFileSaveDir, newFileName)
		localFilePath = fmt.Sprintf("%s/%s", path, newFileName)
	}
	defer file.Close()
	err = ctx.SaveToFile("file", completeFilePath)

	if err != nil {
		log.Fatal(err)
		response.Code = vo.UploadFileFailed
		response.Message = "存储文件失败"
	} else {
		newFileModel, err := models.NewFileAndSave(uploadFilePath, localFilePath, tokenModel.Client)
		if err != nil {
			response.Code = vo.UploadFileFailed
			response.Message = "存储文件记录到数据库错误"
		} else {
			response.Data = make(map[string] interface{})
			response.Data["slug"] = newFileModel.Slug
		}
	}
	ctx.ServeJSON()
}