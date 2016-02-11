package main

import (
	_ "github.com/eraserxp/coedit/routers"
	"github.com/astaxie/beego"
)

func main() {
	//serve static files
	beego.BConfig.WebConfig.StaticDir["/static"] = "static"
	beego.Run()
}

