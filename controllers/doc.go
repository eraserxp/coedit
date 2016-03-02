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

	document_id := c.Ctx.Input.Param(":uuid")

	//if the url doesn't contain uuid, then it means creating a new document
	if (document_id == "") {
		document, err := createDoc()
		if (err != nil) {
			fmt.Println("Failed to create a new document!")
		}
		document_id = document.Id
		c.Redirect("/doc/" + document_id, 302)
	}

	//	c.Data["Website"] = "beego.me"
	//	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "doc.tpl"
}

