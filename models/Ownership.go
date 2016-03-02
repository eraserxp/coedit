package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
)

type Ownership struct {
	Id int
	Username string `orm:"size(20);unique"`
	DocumentId string `orm:"size(36);unique"`
}

