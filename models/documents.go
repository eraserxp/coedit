package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Documents struct {
	Id string `orm:"pk;size(200)"`
	Content string `orm:"type(text)"`
	Privacy string `orm:"size(1)"`
	AccessEmails string `orm:"type(text)"`
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

func (doc *Documents) CheckPrivacyInfo() string {
	var lists []orm.ParamsList

	num, _ := o.Raw(" select privacy from documents where id = ?", doc.Id).ValuesList( &lists)

	if num == 1 {
		return lists[0][0].(string);
	}

	return "D";
}

func (doc *Documents) SearchPrivacyInfo() string {
	var lists []orm.ParamsList
	var result string = ""

	num, _ := o.Raw(" select privacy, access_emails from documents where id = ?", doc.Id).ValuesList( &lists)

	if num == 1 {
		result = result + "{\"privacy\":\"" + lists[0][0].(string) +"\",\"accessemails\":\"" + lists[0][1].(string) +"\"}"
	}


	return result;
}

func (doc *Documents) UpdatePrivacyInfo() error {

	p, _ := o.Raw("update documents set privacy = ? where id = ?").Prepare()

	if _, rerr := p.Exec( doc.Privacy, doc.Id ); rerr != nil {
		if rerr.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", rerr)
			return rerr
		}
	}

	p.Close()

	m, _ := o.Raw("update documents set access_emails = ? where id = ?").Prepare()

	if _, rerr := m.Exec( doc.AccessEmails, doc.Id ); rerr != nil {
		if rerr.Error() != "no LastInsertId available" {
			fmt.Printf("ERR: %v\n", rerr)
			return rerr
		}
	}

	m.Close()

	return nil;
}