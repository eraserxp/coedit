package controllers

import (
	_ "fmt"
	_ "github.com/satori/go.uuid"
	"github.com/satori/go.uuid"
	"github.com/eraserxp/coedit/models"
	"fmt"
	"net/http"
	"encoding/json"
)


func createDoc() (*models.Documents, error) {
	//create document id
	document_id := uuid.NewV4().String()

	//create the document in database
	document := &models.Documents{document_id, ""}
	err := document.Save()
	if (err != nil) {
		fmt.Println("ERROR: " + err.Error())
		return nil, err
	}

	//create token
//	write_token := uuid.NewV4().String()
//	read_token := uuid.NewV4().String()
//	token := &models.Token{document_id, write_token, read_token}

	//save the tokens to postgres
//	err = token.Save()
//	if (err != nil) {
//		fmt.Println("ERROR: " + err.Error())
//		return err
//	}

	fmt.Println("A new document has been created")
	return document, nil
}

type UserNewDocHandler struct {

}

func (this *UserNewDocHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	UserNewDoc(res, req)
}

type newDocReq_struct struct {
	DocumentName string
}

func UserNewDoc (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case "POST":
			sess, _ := globalSessions.SessionStart(res, req)

			UserName := sess.Get("username")

			decoder := json.NewDecoder( req.Body)
			var doc newDocReq_struct
			err := decoder.Decode(&doc)
			if err != nil {
				fmt.Println("Decode Error!")
			}
			fmt.Println( "Create Document Request from UserName: " + UserName.(string) + "; DocumentName: " + doc.DocumentName)

			document, derr := createDoc()

			if derr == nil {
				document_id := document.Id
				fmt.Println("Document ID: " + document_id);

				setExpiredTime(document_id)


				os := &models.Ownership{ 1 , UserName.(string), doc.DocumentName, document_id}
				os.SaveExceptID()
			}

		default:

	}
}