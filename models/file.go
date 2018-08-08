package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"go-file-store/utils"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"os"
	"path"
	"mime"
	"log"
)

type File struct {
	Id uint `orm:"auto;pk"`
	Slug string `orm:"size(32)"`
	Client *Client	`orm:"rel(fk)"` // 一对多外键
	Upload string `orm:"size(256)"`
	Local string `orm:"size(256)"`
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

func (f *File) TableName() string  {
	return "files"
}

func (f *File) TableUnique() [][]string  {
	return [][]string{
		[]string{"Slug"},
	}
}

// NewFileAndSave  创建新文件并保存
func NewFileAndSave(uploadFileName string, client *Client) (*File, error) {
	var slug string
	db := orm.NewOrm()
	file := &File{
		Client: client,
		Upload: uploadFileName,
	}
	for {
		slug = utils.MD5(utils.GenerateRandomString(24))
		file.Slug = slug
		count, err := db.QueryTable(file.TableName()).Filter("slug", slug).Count()
		if err != nil {
			return nil, err
		}
		if count == 0 {
			_, err := db.Insert(file)
			if err != nil {
				return nil, err
			}
			return file, nil
		}
	}
}

// GetFileBySlug 根据Slug获取文件
func GetFileBySlug(slug string) (*File, error)  {
	db := orm.NewOrm()
	file := &File{
		Slug: slug,
	}
	err := db.Read(file, "Slug")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	db.Read(file.Client)
	return file, nil
}


// GetCompleteFilePath 获取文件的绝对路径
func (f *File) GetCompleteFilePath() string {
	return fmt.Sprintf("%s/%s/%s",
		beego.AppConfig.String("uploadFilesDir"),
		f.Client.ClientId,
		strings.TrimLeft(f.Local, "/"))
}

// isFileExists 判断文件是否存在
func (f *File) IsFileExists() bool {
	if _, err := os.Stat(f.Upload); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// FileName	文件名
func (f *File) FileName() string {
	segments := strings.Split(f.Upload, "/")
	return segments[len(segments) - 1]
}

// FileExt 返回文件的扩展名
func (f *File) FileExt() string {
	fileName := f.FileName()
	return path.Ext(fileName)
}

// FileMIME	根据扩展名返回文件类型
func (f *File) FileMIME() string {
	return mime.TypeByExtension(f.FileExt())
}