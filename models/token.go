package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
)

type Token struct {
	DocumentId string `orm:"pk;size(36)"`
	WriteToken string `orm:"size(36);unique"`
	ReadToken string `orm:"size(36);unique"`
}


func (token *Token) Save() error {
	if _, err := o.Insert(token); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}
