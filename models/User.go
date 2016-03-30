package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
)

type User struct {
	Name string `orm:"pk;size(20)"`
	Password string `orm:"size(20)"`
}

func (u *User) SaveUser() error {
	if _, err := o.Insert(u); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}


