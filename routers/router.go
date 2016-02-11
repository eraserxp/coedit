package routers

import (
	"github.com/eraserxp/coedit/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
