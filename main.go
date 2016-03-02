package main

import (
	"fmt"
	_ "github.com/eraserxp/coedit/routers"
	"github.com/astaxie/beego"
)


func main() {
	//serve static files
	beego.BConfig.WebConfig.StaticDir["/doc/static"] = "static"
	fmt.Println("write into database")

	beego.Run()
}

