package models

import (
	"github.com/astaxie/beego"
	"fmt"
	"strings"
	"github.com/astaxie/beego/orm"
)

// 初始化数据库连接
func init() {
	dbHost := beego.AppConfig.String("dbHost")
	dbPort := beego.AppConfig.String("dbPort")
	dbUser := beego.AppConfig.String("dbUser")
	dbPassword := beego.AppConfig.String("dbPassword")
	dbDatabase := beego.AppConfig.String("dbDatabase")

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	dbSource = strings.Replace(dbSource, "Local", "Asia%2FShanghai", 1)
	orm.RegisterDataBase("default", "mysql", dbSource, 30, 30)
}
