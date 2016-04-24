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
		c.TplName="404.tpl"
		return;
	} else if ( privacyOption == "S") {
		c.Redirect("/docreg/" + document_id, 302)
	}

	//	c.Data["Website"] = "beego.me"
	//	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "doc.tpl"
}

type DocRegController struct {
	beego.Controller
}

func (c *DocRegController) Get() {
	fmt.Println("DocRegController")

	document_id := c.Ctx.Input.Param(":uuid")

	if ( document_id == "") {

		c.TplName = "404.tpl"

	} else {

		sess, _ := globalSessions.SessionStart( c.Ctx.ResponseWriter, c.Ctx.Request)
		username := sess.Get("username")



		if( username == nil || username == "") {
			c.TplName = "404.tpl"
			fmt.Println( "Unregistered user attempting to load registered page." );
			return;
		}

		fmt.Println( username.(string) + " is attempting to read file " + document_id);

		doc := &models.Documents{document_id, "", "E", ""}

		canAccess := doc.CheckAccessible( username.(string) )


		c.Data["Email"] =  username.(string)

		if ( !canAccess ) {
			c.TplName = "404.tpl"
			fmt.Println(username.(string) + " can not access the page" )
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



		c.TplName = "regdoc.tpl"
	}
}

type RegDocController struct {
	beego.Controller
}

func (c *RegDocController) Get() {
	fmt.Println("RegDocController")

	document_id := c.Ctx.Input.Param(":uuid")

	if ( document_id == "") {

		c.TplName = "404.tpl"

	} else {

		sess, _ := globalSessions.SessionStart( c.Ctx.ResponseWriter, c.Ctx.Request)
		username := sess.Get("username")

		if( username == nil || username == "") {
			c.TplName = "404.tpl"
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



		c.TplName = "regdoc.tpl"
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
