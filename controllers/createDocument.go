package controllers

import (
	_ "fmt"
	_ "github.com/satori/go.uuid"
	"github.com/satori/go.uuid"
	"github.com/eraserxp/coedit/models"
	"fmt"
	"net/http"
	"encoding/json"
	"io"
)


func createDoc() (*models.Documents, error) {
	//create document id
	document_id := uuid.NewV4().String()

	//create the document in database
	document := &models.Documents{document_id, "", "E", ""}
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

			if( doc.DocumentName == "") {
				io.WriteString(res, "EMPTY")
				return;
			}

			searchName := &models.Ownership{ 1 , UserName.(string), doc.DocumentName, "default"}
			dupCheckResult := searchName.SearchDupName();

			if ( !dupCheckResult ) {

				io.WriteString(res, "Dup")

			} else {

				document, derr := createDoc()

				if derr == nil {
					document_id := document.Id
					fmt.Println("Document ID: " + document_id);

					// user document should not have expire time
					//setExpiredTime(document_id)

					os := &models.Ownership{1, UserName.(string), doc.DocumentName, document_id}

					os.SaveExceptID()
				}

				io.WriteString(res, "OK")
			}

		default:

	}
}

type DeleteDocHandler struct {

}

func (this *DeleteDocHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	DeleteDoc(res, req)
}

type deleteDoc_struct struct {
	DocumentName string
}

func DeleteDoc (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		sess, _ := globalSessions.SessionStart(res, req)

		UserName := sess.Get("username")

		decoder := json.NewDecoder( req.Body)
		var doc deleteDoc_struct
		err := decoder.Decode(&doc)
		if err != nil {
			fmt.Println("Decode Error!")
		}
		fmt.Println( "Create Document Request from UserName: " + UserName.(string) + "; DocumentName: " + doc.DocumentName)

		os := &models.Ownership{ 1 , UserName.(string), doc.DocumentName, "default"}

		doc_id := os.SearchID();

		doc_instance := &models.Documents{ doc_id, "default", "E", ""}

		doc_instance.Delete();

	default:

	}
}