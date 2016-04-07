package models

import (
	_ "fmt"
	_ "github.com/astaxie/beego/orm"
	_ "github.com/lib/pq" // import postgres driver
	"fmt"
	"github.com/astaxie/beego/orm"
	"encoding/json"
)

type Access struct {
	DocumentId string `orm:"size(36)`
	
}

