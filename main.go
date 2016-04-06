package main

import (
	"fmt"
	_ "github.com/eraserxp/coedit/routers"
	"github.com/astaxie/beego"
	"net/http"
	"net/url"
	"github.com/eraserxp/coedit/websocketproxy"
	"log"
)

func startWebsocketProxy() error {
	u, _ := url.Parse("ws://localhost:8001")
	err := http.ListenAndServe(":8000", websocketproxy.NewProxy(u))
	if err != nil {
		log.Fatalln(err)
		return err
	}
	fmt.Println("start the websocket proxy")
	return nil
}

func hook() error  {
	go startWebsocketProxy()
	return nil
}

//func init()  {
//	startWebsocketProxy()
//}

func main() {
	//run the websocket proxy
	//startWebsocketProxy()
	beego.AddAPPStartHook(hook)

	//run the web server to serve the static files
	beego.BConfig.WebConfig.StaticDir["/doc/static"] = "static"
	fmt.Println("write into database")

	//set session

	beego.Run()
}




