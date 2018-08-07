package controllers

import (
	"go-file-store/vo"
	"strings"
	"go-file-store/models"
	"path"
	"github.com/astaxie/beego"
	"fmt"
	"os"
	"time"
	"log"
	"strconv"
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

	token := strings.Trim(ctx.GetString("token"), " ")
	tokenModel, err := models.IsTokenValidate(token)

	if err != nil {
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
	filePath, err := saveFiles(ctx, tokenModel.Client, response)
	if err != nil {
		ctx.ServeJSON()
		ctx.StopRun()
	}
	newFileModel, err := models.NewFileAndSave(filePath, tokenModel.Client)
	if err != nil {
		response.Code = vo.UploadFileFailed
		response.Message = "存储文件记录到数据库错误"
	} else {
		response.Data = make(map[string] interface{})
		response.Data["slug"] = newFileModel.Slug
	}
	ctx.ServeJSON()
}

// saveFiles 保存文件方法
func saveFiles(ctx *FileController, client *models.Client, response *vo.ResponseMessage) (string, error)  {
	file, fileHeader, err := ctx.GetFile("file")
	if err != nil {
		response.Code = vo.UploadFileFailed
		response.Message = "读取文件错误"
		return "", err
	}
	targetSaveFolder := client.RootDir()
	if _, err = os.Stat(targetSaveFolder); err != nil && os.IsNotExist(err) {
		os.MkdirAll(targetSaveFolder, os.ModePerm)
	}
	fileExt := path.Ext(fileHeader.Filename)
	uploadFilePath := fmt.Sprintf("%s/%s%s", targetSaveFolder, strconv.FormatInt(time.Now().UnixNano(), 10), fileExt)
	defer file.Close()
	err = ctx.SaveToFile("file", uploadFilePath)
	if err != nil {
		log.Fatal(err)
		response.Code = vo.UploadFileFailed
		response.Message = "存储文件失败"
		return "", err
	}
	return uploadFilePath, nil
}