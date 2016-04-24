package tests

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/eraserxp/coedit/models"
	"github.com/astaxie/beego/orm"
	"fmt"
)

func TestCreateDoc(t *testing.T) {
	var o  orm.Ormer = orm.NewOrm()
	docNew := &models.Documents{"testid","test_content","D","test@gmail"}
	// Only pass t into top-level Convey calls
	err := docNew.Save()
	if (err == nil) {
		fmt.Println("***************************************")
	}
	Convey("Test Create Doc in database", t, func() {
		var lists []orm.ParamsList

		num, _ := o.Raw(" select Id from documents where id = ?", "testid").ValuesList( &lists)

		num=num+1
		So(lists[0][0], ShouldEqual, "testid")
	})
}
