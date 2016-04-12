package routers

import (
	"github.com/eraserxp/coedit/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Handler("/auth/:provider", &controllers.AuthHandler{})
	beego.Handler("/auth/:provider/callback", &controllers.AuthCallbackHandler{})
	beego.Router("/doc/?:uuid", &controllers.DocController{})
	beego.Router("/regdoc/?:uuid", &controllers.RegDocController{})
	beego.Router("/user/?:uemail", &controllers.AccountController{})
	beego.Router("/", &controllers.MainController{})


	beego.Handler("/addnewdoc", &controllers.UserNewDocHandler{})
	beego.Handler("/requestuserlist", &controllers.RequestUserListHandler{})
	beego.Handler("/opendoc", &controllers.OpenDocReqHandler{})
	beego.Handler("/logout", &controllers.LogoutHandler{} )
	beego.Handler("/deletedoc", &controllers.DeleteDocHandler{})
}
