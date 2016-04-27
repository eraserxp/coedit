package controllers

import (
	"github.com/astaxie/beego"
	"github.com/eraserxp/coedit/models"
	"fmt"
	"net/http"
	"encoding/json"
	"io"
)

type DocController struct {
	beego.Controller
}

func (c *DocController) Get() {
	fmt.Println("DocController")

	document_id := c.Ctx.Input.Param(":uuid")

	//set document id as cookies
	//cookie := http.Cookie{Name: "documentId", Value: document_id}
	//http.SetCookie(c.Ctx.ResponseWriter, &cookie)

	sess, _ := globalSessions.SessionStart( c.Ctx.ResponseWriter, c.Ctx.Request)

	sess.Set("documentId", document_id);

	//if the url doesn't contain uuid, then it means creating a new document
	if (document_id == "") {
		document, err := createDoc()
		if (err != nil) {
			fmt.Println("Failed to create a new document!")
		}
		document_id = document.Id
		//set the expire time for the document
		setExpiredTime(document_id)
		c.Redirect("/doc/" + document_id, 302)
	}

	docmodel := &models.Documents{document_id, "", "E", ""};
	privacyOption := docmodel.CheckPrivacyInfo()

	if( privacyOption == "N" ) {
		c.TplName="noAccess.tpl"
		return;
	} else if ( privacyOption == "S") {
		c.Redirect("/docreg/" + document_id, 302)
	}

	//	c.Data["Website"] = "beego.me"
	//	c.Data["Email"] = "astaxie@gmail.com"
	sess.Set("previousDoc", "none")

	c.TplName = "doc.tpl"
}

//controller for other users to access documents created by registered users
type DocRegController struct {
	beego.Controller
}

func (c *DocRegController) Get() {
	fmt.Println("DocRegController")

	document_id := c.Ctx.Input.Param(":uuid")

	if ( document_id == "") {
		c.TplName = "404.tpl"
	} else {
		//set document id as cookies
		//cookie := http.Cookie{Name: "documentId", Value: document_id}
		//http.SetCookie(c.Ctx.ResponseWriter, &cookie)

		sess, _ := globalSessions.SessionStart( c.Ctx.ResponseWriter, c.Ctx.Request)
		username := sess.Get("username")

		sess.Set("documentId", document_id);

		//check if the document can be accessed by everyone
		docmodel := &models.Documents{document_id, "", "E", ""};
		privacyOption := docmodel.CheckPrivacyInfo()
		//use doc.tpl if user is not logged
		if privacyOption == "E" && (username == nil || username == "") {
			sess.Set("previousDoc", "none")
			c.TplName = "doc.tpl"
			return
		}

		if( username == nil || username == "") {
			c.TplName = "noAccess.tpl"
			fmt.Println( "Unregistered user attempting to load registered page." );
			return;
		}

		fmt.Println( username.(string) + " is attempting to read file " + document_id);

		doc := &models.Documents{document_id, "", "E", ""}

		canAccess := doc.CheckAccessible( username.(string) )


		c.Data["Email"] =  username.(string)

		if ( !canAccess ) {
			fmt.Println(username.(string) + " can not access the page" )
			c.TplName = "noAccess.tpl"
			return;
		} else {
			fmt.Println(username.(string) + " can access the page" )
		}


		filename := ( &models.Ownership{1, username.(string), "default", document_id} ).SearchDocName()

		if( filename == "") {
			c.TplName = "404.tpl"
			return;
		}

		c.Data["FileName"] =  filename


		sess.Set("previousDoc", "none")

		c.TplName = "docreg.tpl"
	}
}

// controller for an owner to access his own documents
type RegDocController struct {
	beego.Controller
}


func (c *RegDocController) Get() {
	fmt.Println("RegDocController")

	document_id := c.Ctx.Input.Param(":uuid")

	if ( document_id == "") {

		c.TplName = "noAccess.tpl"

	} else {
		//set document id as cookies
		//cookie := http.Cookie{Name: "documentId", Value: document_id}
		//http.SetCookie(c.Ctx.ResponseWriter, &cookie)

		sess, _ := globalSessions.SessionStart( c.Ctx.ResponseWriter, c.Ctx.Request)
		username := sess.Get("username")

		sess.Set("previousDoc", document_id)

		sess.Set("documentId", document_id);

		//check if privacy option of the document
		docmodel := &models.Documents{document_id, "", "E", ""};
		privacyOption := docmodel.CheckPrivacyInfo()

		if privacyOption == "N" { //no other users can access
			c.TplName = "noAccess.tpl"
			return
		}

		//everyone can access and the current user is not logged in
		if privacyOption == "E" && (username == nil || username == "") {
			sess.Set("previousDoc", "none")
			c.TplName = "doc.tpl"
			return
		}

		if( username == nil || username == "") {
			//c.TplName = "404.tpl"
			c.TplName = "login.tpl"
			fmt.Println( "Unregistered user attempting to load registered page." );
			return;
		}

		c.Data["Email"] =  username.(string)
		filename := ( &models.Ownership{1, username.(string), "default", document_id} ).SearchDocName()

		if( filename == "") {
			c.TplName = "404.tpl"
			return;
		}

		c.Data["FileName"] =  filename

		owner := ( &models.Ownership{1, "", "default", document_id} ).GetOwner()

		sess.Set("previousDoc", "none")

		if owner != username {
			c.TplName = "docreg.tpl" //no way to change the privacy option of a document
		} else {
			c.TplName = "regdoc.tpl" //can change the privacy option of a document
		}

	}
}

type OpenDocReqHandler struct {

}

func (this *OpenDocReqHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	OpenDocReq(res, req)
}

type OpenDocReq_struct struct {
	DocumentName string
}

func OpenDocReq (res http.ResponseWriter, req *http.Request) {
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
		fmt.Println( "Open Document Request from UserName: " + UserName.(string) + "; DocumentName: " + doc.DocumentName)

		os := &models.Ownership{ 1 , UserName.(string), doc.DocumentName, "default"}
		docID := os.SearchID()

		fmt.Println("Return doc ID: " + docID)
		io.WriteString( res, docID)

	default:

	}
}

type LoadDocPrivacyHandler struct {

}

func (this *LoadDocPrivacyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	LoadDocPrivacy(res, req)
}

type LoadDocPrivacy_struct struct {
	DocumentName string
}

func LoadDocPrivacy (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		sess, _ := globalSessions.SessionStart(res, req)

		UserName := sess.Get("username")

		decoder := json.NewDecoder( req.Body)
		var doc LoadDocPrivacy_struct
		err := decoder.Decode(&doc)
		if err != nil {
			fmt.Println("Decode Error!")
		}
		fmt.Println( "Load Document Privacy Request from UserName: " + UserName.(string) + "; DocumentName: " + doc.DocumentName)

		os := &models.Ownership{ 1 , UserName.(string), doc.DocumentName, "default"}
		docID := os.SearchID()

		fmt.Println("Return doc ID: " + docID)
		docmodel := &models.Documents{docID, "", "E", ""}

		msg := docmodel.SearchPrivacyInfo();

		fmt.Println(msg);

		io.WriteString( res, msg)

	default:

	}

}


type SaveDocPrivacyHandler struct {

}

func (this *SaveDocPrivacyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	SaveDocPrivacy(res, req)
}

type SaveDocPrivacy_struct struct {
	DocumentName string
	Privacy string
	AccessEmails string
}

func SaveDocPrivacy (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		sess, _ := globalSessions.SessionStart(res, req)

		UserName := sess.Get("username")

		decoder := json.NewDecoder( req.Body)
		var doc SaveDocPrivacy_struct
		err := decoder.Decode(&doc)
		if err != nil {
			fmt.Println("Decode Error!")
		}
		fmt.Println( "Save Document Privacy Request from UserName: " + UserName.(string) + "; DocumentName: " + doc.DocumentName)

		os := &models.Ownership{ 1 , UserName.(string), doc.DocumentName, "default"}
		docID := os.SearchID()

		fmt.Println("Return doc ID: " + docID)
		docmodel := &models.Documents{docID, "", doc.Privacy, doc.AccessEmails}

		updateerr := docmodel.UpdatePrivacyInfo();

		if updateerr == nil {
			io.WriteString( res, "OK")
		} else {
			io.WriteString( res, "ERROR")
		}

	default:

	}

}
