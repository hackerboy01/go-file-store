// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"go-file-store/controllers"
)

// 初始化路由配置
func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSBefore(func(ctx *context.Context) {
			ctx.Output.Header("Access-Control-Allow-Origin", ctx.Input.Domain())
			ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
			ctx.Output.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		}),
		beego.NSRouter("client", &controllers.ClientController{}),
		beego.NSRouter("token", &controllers.TokenController{}),
		beego.NSAutoRouter(&controllers.FileController{}),
	)
	beego.AddNamespace(ns)
}
