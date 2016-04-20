package tests


import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	_ "fmt"
	_ "github.com/lib/pq" // import postgres driver
	"github.com/eraserxp/coedit/models"
)

func TestOwnership(t *testing.T) {
	ownTest := &models.Ownership{666, "test1", "testFile", "abc123DEF"};

	Convey("Check save ownership", t, func() {

		success := ownTest.Save()
		So(success, ShouldEqual, nil)
	})

	Convey("Check ownership document id", t, func() {

		So(ownTest.DocumentId, ShouldEqual, ownTest.SearchID())
	})

	Convey("Check ownership filename ", t, func() {

		So(ownTest.Filename, ShouldEqual, ownTest.SearchDocName())
	})
}