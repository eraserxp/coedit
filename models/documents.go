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
	if _, err := o.Insert(doc); err != nil {
		if err.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", err)
			return err
		}
	}
	return nil
}

func (doc *Documents) Delete() error {
	p, _ := o.Raw("delete from ownership where document_id = ?").Prepare()

	if _, rerr := p.Exec( doc.Id ); rerr != nil {
		if rerr.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", rerr)
			return rerr
		}
	}

	p.Close()

	s, _ := o.Raw("delete from expire where document_id = ?").Prepare()

	if _, rerr := s.Exec( doc.Id ); rerr != nil {
		if rerr.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", rerr)
			return rerr
		}
	}

	s.Close()

	m, _ := o.Raw("delete from documents where id = ?").Prepare()

	if _, rerr := m.Exec( doc.Id ); rerr != nil {
		if rerr.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", rerr)
			return rerr
		}
	}

	m.Close()

	return nil
}