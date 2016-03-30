package routers

import (
	"github.com/eraserxp/coedit/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Handler("/auth/:provider", &controllers.AuthHandler{})
	beego.Handler("/auth/:provider/callback", &controllers.AuthCallbackHandler{})
	beego.Router("/doc/?:uuid", &controllers.DocController{})
	beego.Router("/user/?:uemail", &controllers.AccountController{})
	beego.Router("/", &controllers.MainController{})
}
