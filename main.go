package main

import (
	_ "coedit/routers"
	"github.com/astaxie/beego"
)

func main() {
	//serve static files
	beego.StaticDir["/static"] = "static"
	beego.Run()
}

