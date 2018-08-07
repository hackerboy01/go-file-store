package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"go-file-store/utils"
	"log"
	"fmt"
	"strings"
	"github.com/astaxie/beego"
)

type Client struct {
	Id uint `orm:"auto;pk"`
	ClientId string `orm:"size(128);unique;column(client_id)"`
	ClientSecret string `orm:"size(32);column(client_secret)"`
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

func (client *Client) TableName() string  {
	return "clients"
}

// NewClient 生成新的Client
func NewClient() (*Client, error)  {
	clientId, err := generateClientId()
	log.Print(clientId)
	if err == nil {
		client := &Client{
			ClientId: clientId,
			ClientSecret: generatePassword(12),
		}
		return client, nil
	}

	return nil, err
}

// Save 保存账号
func (client *Client) Save() (id int64, err error)  {
	db := orm.NewOrm()
	client.ClientSecret = utils.CryptPassword(client.ClientSecret)
	id, err = db.Insert(client)
	return id, err
}

// ValidateClient 验证client_id 和 client_secret
func ValidateClient(clientId string, clientSecret string) (*Client, error)  {
	db := orm.NewOrm()
	newClientSecret := utils.CryptPassword(clientSecret)
	client := Client{
		ClientId: clientId,
		ClientSecret: newClientSecret,
	}
	err := db.Read(&client, "ClientId", "ClientSecret")
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (client *Client) RootDir() string {
	clientFolder := fmt.Sprintf("%s/%s",
		strings.TrimRight(beego.AppConfig.String("uploadFilesDir"), "/"),
		client.ClientId)
	return clientFolder
}

// generateClientId 生成账户Id
func generateClientId() (clientId string, err error)  {
	numbers := "0123456789"
	sequences := make([]byte, 16)
	db := orm.NewOrm()
	client := &Client{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		for i := range sequences {
			sequences[i] = numbers[r.Intn(len(numbers))]
		}
		clientId = string(sequences)
		count, err := db.QueryTable(client.TableName()).Filter("client_id", clientId).Count()
		if err != nil {
			return "", err
		}
		if count == 0 {
			return clientId, nil
		}
	}
}

// generatePassword	生成密码
func generatePassword(length uint) string  {
	return utils.GenerateRandomString(length)
}