package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
)

type Ownership struct {
	Id int
	Username string `orm:"size(40)"`
	Filename string `orm:"size(100)"`
	DocumentId string `orm:"size(36);unique"`
}

func (ownership *Ownership) Save() error {
	if _, err := o.Insert(ownership); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}