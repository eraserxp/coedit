package main

import (
	_ "coedit/routers"
	"github.com/astaxie/beego"
)

func main() {
//	beego.SetStaticPath("/static","static")
	beego.StaticDir["/static"] = "/static"
	beego.Run()
}

