package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
)

type Documents struct {
	Id string `orm:"pk;size(200)"`
	Content string `orm:"type(text)"`
}


func (doc *Documents) Save() error {
	fmt.Println("inside token save")
	if _, err := o.Insert(doc); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}

