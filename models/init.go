package models

import (
	_ "fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
)

//global variable for the model package
var o orm.Ormer

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)

	// set default database
	orm.RegisterDataBase("default", "postgres", "user=leaps password=leaps dbname=leaps sslmode=disable", 30)

	// register model
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Token))
	orm.RegisterModel(new(Ownership))
	orm.RegisterModel(new(Documents))
	orm.RegisterModel(new(Expire))


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

	o = orm.NewOrm()

	//insert an none user
	user := User{Name: "none", Password: "none"}
	if _, err := o.Insert(&user); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
		}
	}

}
