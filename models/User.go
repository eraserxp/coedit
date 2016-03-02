package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
)

type User struct {
	Name string `orm:"pk;size(20)"`
	Password string `orm:"size(20)"`
}




