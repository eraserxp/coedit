package main

import (
	_ "coedit/routers"
	"github.com/astaxie/beego"
)

func main() {
	//serve static files
	beego.BConfig.WebConfig.StaticDir["/static"] = "static"
	beego.Run()
}

