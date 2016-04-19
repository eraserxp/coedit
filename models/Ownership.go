package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
	"github.com/astaxie/beego/orm"
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

func (ownership *Ownership) SaveExceptID() error {
	p, _ := o.Raw("insert into ownership(username, filename, document_id) values (?, ?, ?)").Prepare()

	if _, rerr := p.Exec( ownership.Username, ownership.Filename, ownership.DocumentId ); rerr != nil {
		if rerr.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", rerr)
			return rerr
		}
	}

	p.Close()
	return nil
}

func (ownership * Ownership) SearchDupName() bool {
	var lists []orm.ParamsList

	num, _ := o.Raw(" select document_id from ownership where username = ? and filename = ?", ownership.Username, ownership.Filename).ValuesList( &lists)

	if( num == 0) {
		return true;
	} else {
		return false;
	}
}

func (ownership *Ownership) SearchID() string {

	var lists []orm.ParamsList

	num, err := o.Raw(" select document_id from ownership where username = ? and filename = ?", ownership.Username, ownership.Filename).ValuesList( &lists)

	if err == nil {
		if num == 1 {
			return lists[0][0].(string)

		} else {
			fmt.Println("No result found or result number is not correct !")
		}
	} else {
		fmt.Println( "Error on select ownership query! %v" , err)
	}

	return ""
}

func (ownership *Ownership) SearchDocName() string {

	var lists []orm.ParamsList

	num, err := o.Raw(" select filename from ownership where document_id = ?", ownership.DocumentId).ValuesList( &lists)

	if err == nil {
		if num == 1 {
			return lists[0][0].(string)

		} else {
			fmt.Println("No result found or result number is not correct !")
		}
	} else {
		fmt.Println( "Error on select ownership query! %v" , err)
	}

	return ""
}