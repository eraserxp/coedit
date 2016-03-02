package models

import (
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
	"time"
)

type Expire struct {
	DocumentId string `orm:"pk;size(36)"`
	ExpireTime time.Time `orm:"datetime"`
}


func (expire *Expire) Save() error {
	if _, err := o.Insert(expire); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}


func (expire *Expire) Update(expireTime time.Time) error {
	if _, err := o.Insert(expire); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}