package routers

import (
	"github.com/eraserxp/coedit/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/doc", &controllers.DocController{})
    beego.Router("/", &controllers.MainController{})
}
