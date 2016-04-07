package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
	"github.com/astaxie/beego/orm"
	"encoding/json"
)

type Account struct {
	Name string `orm:"pk;size(40)"`
	Password string `orm:"size(20)"`
}

/*func (a *Account) SaveAcc() error {
	if _, err := o.Insert(a); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}*/

func (a *Account) CheckExist() bool {

	if err:= o.Read( a, "Name") ; err == nil {
		fmt.Println("Found a Record for " + a.Name)
		return false
	} else {
		fmt.Println("Insert a new record for " + a.Name)
		o.Insert(a)
		return true
	}
}


func (a *Account) SearchDocument() string {
	var lists []orm.ParamsList
	result := ""

	num, err := o.Raw("SELECT filename FROM ownership where username = ?", a.Name).ValuesList( &lists)

	if err == nil {
		if num == 0 {
			fmt.Println("No result found")
		} else {
			fmt.Println( "Found result!" )

			var docList = make([]string, num)

			for i := 0; i < int(num); i++ {
				docList[i], _ = lists[i][0].(string)
			}

			r, _ := json.Marshal( docList)
			result = string(r)
		}

	} else {
		fmt.Println(err)

	}

	return result
}
