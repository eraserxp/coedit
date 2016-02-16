package test

import (
	"github.com/eraserxp/coedit/controllers"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)



func TestConnectRedis(t *testing.T) {

	key1 := "yao test"
	value := "dfasdfa"

	success := controllers.PutToRedis(key1, value)

	if (success == 1){
		fmt.Println("put success")
	} else {
		fmt.Println("put failed")
	}


	getV := controllers.GetFromRedis(key1)

//	if (getV == value){
//		fmt.Println("get value and matched")
//	} else {
//		fmt.Println("get wrong value")
//	}
	Convey("The get reuslt match the put value", t, func() {
		So( getV, ShouldEqual, value)
	})

}


