package tests

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/eraserxp/coedit/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"time"
	"fmt"
)

//global variable for the model package
var o orm.Ormer = orm.NewOrm()

func GetExpiredTime(DocumentId string) string {
	var lists []orm.ParamsList
	//fmt.Println( "My document Id: %s" , DocumentId)
	num, err := o.Raw("SELECT expire_time FROM expire WHERE document_id = ?", DocumentId).ValuesList( &lists)
	//fmt.Println( "My expire time: %s" , expireTime)
	if err == nil {
		if (num == 1){
			return lists[0][0].(string)
		} else {
			fmt.Println( "Error on select expire query! %v" , err)
		}
	} else {
		fmt.Println( "Error on select expire query! %v" , err)
	}

	return ""
}

func timeToDT(t time.Time) string {
	var lists []orm.ParamsList
	e := models.Expire{"TestDID", t}
	e.Save()
	o.Raw("SELECT expire_time FROM expire WHERE document_id = ?", "TestDID").ValuesList( &lists)
	o.Raw("DELETE FROM expire WHERE document_id = ?", "TestDID").Exec()
	return lists[0][0].(string)
}

func isSameTime(a string, b string) int {
	// this is not 100% right method to compare the two time strs
	// but should be enough for now
	// should be improved in future
	count := 0
	for i := 0; i < len(a) && i < len(b); i++ {
		if (a[i] == b[i]){
			count++
		}
	}
	return count;
}

func TestSpec(t *testing.T) {
	Convey("Given a expire time and save it to the database", t, func() {
		curTime := time.Now()
		exp := models.Expire{"76441e7c-310c-405a-89ce-c39bcd288c03", curTime}
		exp.Save()
		temp := GetExpiredTime(exp.DocumentId)
		Convey("Returned expire time should not be empty", func() {
			So(temp, ShouldNotBeEmpty)
		})
		Convey("Expire time should be correctly saved", func() {
			count := isSameTime(temp, timeToDT(exp.ExpireTime));
			So(count, ShouldBeGreaterThan, 19)
		})
		Convey("Expire time should be correctly updated", func() {
			exp.Update(time.Now())
			uTemp := GetExpiredTime(exp.DocumentId)
			count := isSameTime(uTemp, timeToDT(exp.ExpireTime));
			So(count, ShouldBeGreaterThan, 19)
		})
	})
}
