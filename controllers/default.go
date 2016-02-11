package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	fmt.Print("123")
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index2.tpl"

}
