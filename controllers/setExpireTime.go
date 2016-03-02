package controllers

import (
	"github.com/eraserxp/coedit/models"
	"fmt"
	"time"
)


func setExpiredTime(document_id string) (error) {
	now := time.Now()
	oneDay, _:= time.ParseDuration("24h")
	expired_time := now.Add(oneDay)
	expire := &models.Expire{document_id, expired_time}
	expire.Save()
	fmt.Println("set the expire time for document: " + document_id)
	return nil
}

