package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"go-file-store/utils"
)

type Token struct {
	Id uint `orm:"auto;pk"`
	Token string `orm:"size(32)"`
	Client *Client `orm:"rel(fk)"` // 设置一对多关系
	Expires uint
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

// 自定义表名
func (t *Token) TableName() string  {
	return "tokens"
}

// 自定义唯一索引
func (t *Token) TableUnique() [][]string  {
	return [][]string{
		[]string{"Token"},
	}
}

// NewTokenAndSave 生成Token并存储
func NewTokenAndSave(seconds uint, client *Client) (t *Token, err error)  {
	db := orm.NewOrm()
	token := &Token{}
	for {
		tokenString := utils.GenerateToken()
		token.Token = tokenString
		token.Client = client
		token.Expires = uint(time.Now().Unix() + int64(seconds))
		count, err := db.QueryTable(token.TableName()).Filter("token", tokenString).Count()
		if err != nil {
			return nil, err
		}
		if count == 0 {
			_, err := db.Insert(token)
			if err != nil {
				return nil, err
			}
			return token, nil
		}
	}
}

// IsTokenExist 验证Token是否存在
func IsTokenValidate(token string) (*Token, error)  {
	db := orm.NewOrm()
	tokenModel := &Token{
		Token: token,
	}

	err := db.Read(tokenModel, "Token")
	if err != nil {
		return nil, err
	}
	db.Read(tokenModel.Client)
	return tokenModel, nil
}

// IsTokenExpire验证Token是否过期
func (t *Token) IsTokenExpire() bool {
	if t.Expires < uint(time.Now().Unix()) {
		return true
	}
	return false
}