package main

import (
	"fmt"
	_ "github.com/eraserxp/coedit/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/eraserxp/coedit/models"
)

func setupDatabase()  {
	//set up the database
	// Database alias.
	name := "default"

	// Drop table and re-create.
	//	force := false
	force := true

	// Print log.
	verbose := true

	// Error.
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	o := orm.NewOrm()
	user := models.User{Name: "none", Password: "none"}
	if _, err := o.Insert(&user); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
		}
	}

}

func main() {
	setupDatabase()

	//serve static files
	beego.BConfig.WebConfig.StaticDir["/static"] = "static"
	fmt.Println("write into database")

	beego.Run()
}

