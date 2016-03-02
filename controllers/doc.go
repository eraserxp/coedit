package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
)

type DocController struct {
	beego.Controller
}

func (c *DocController) Get() {
	fmt.Println("DocController")
	createDoc()
	//	c.Data["Website"] = "beego.me"
	//	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "doc.tpl"
}

