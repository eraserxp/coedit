package controllers

import (
	_ "fmt"
	_ "github.com/satori/go.uuid"
	"github.com/satori/go.uuid"
	"github.com/eraserxp/coedit/models"
	"fmt"
)


func createDoc() error {
	//create document id and tokens
	document_id := uuid.NewV4().String()
	write_token := uuid.NewV4().String()
	read_token := uuid.NewV4().String()
	token := &models.Token{document_id, write_token, read_token}

	//save the tokens to postgres
	err := token.Save()
	if (err != nil) {
		fmt.Println("ERROR: " + err.Error())
		return err
	}

	fmt.Println("A new document has been created")
	return nil
}