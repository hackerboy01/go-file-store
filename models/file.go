package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"go-file-store/utils"
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
func NewFileAndSave(uploadFileName string, loaclFileName string, client *Client) (*File, error) {
	var slug string
	db := orm.NewOrm()
	file := &File{
		Client: client,
		Upload: uploadFileName,
		Local: loaclFileName,
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

