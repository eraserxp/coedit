package models

import (
	_ "fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
)


func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)

	// set default database
	orm.RegisterDataBase("default", "postgres", "user=leaps password=leaps dbname=leaps sslmode=disable", 30)

	// register model
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Token))
	orm.RegisterModel(new(Ownership))

}
