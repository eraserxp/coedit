package tests


import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	_ "fmt"
	_ "github.com/lib/pq" // import postgres driver
	"github.com/eraserxp/coedit/models"
	"github.com/astaxie/beego/orm"
)




func TestToken(t *testing.T) {
	tokenTest := &models.Token{"token1", "writeToken", "readToken"};
	var o orm.Ormer = orm.NewOrm()
	result := []orm.ParamsList{}

	Convey("Check save token", t, func() {
		success := tokenTest.Save()
		So(success, ShouldEqual, nil)
	})

	Convey("Check write token", t, func() {

		r, _ := o.Raw("select write_token from token where document_id = ?", tokenTest.DocumentId).ValuesList(&result);
		_ = r
		e := result[0][0]
			So(e, ShouldEqual, tokenTest.WriteToken)
	})

	Convey("Check read token", t, func() {

		r, _ := o.Raw("select read_token from token where document_id = ?", tokenTest.DocumentId).ValuesList(&result);
		_ = r
		e := result[0][0]
			So(e, ShouldEqual, tokenTest.ReadToken)
	})

}

