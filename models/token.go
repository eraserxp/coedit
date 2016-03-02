package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
)

type Token struct {
	DocumentId string `orm:"pk;size(36)"`
	WriteToken string `orm:"size(36);unique"`
	ReadToken string `orm:"size(36);unique"`
}



