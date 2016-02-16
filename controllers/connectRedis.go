package controllers

import (
	"fmt"
	// import redigo lib
	"github.com/garyburd/redigo/redis"
	"time"
)

//var db int
var rs redis.Conn
var connectErr error


func connectRedis() (redis.Conn, error){

	rs , err := redis.DialTimeout("tcp", ":6379", 0, 1*time.Second, 1*time.Second)
	// If error occurs, print err messages and return
	if err != nil {
		fmt.Println(err)
		fmt.Println("redis connect error")
		return nil, nil
	}

	// choose db
	//rs.Do("SELECT", db)

	return rs, err
}


/**
To put key value into redis, if return 1, means success otherwise, fails
 */
func PutToRedis(key string, value string) int{
	rs, connectErr = connectRedis()
	n, err := rs.Do("SET", key, value)
	_ = n
	// if error, return
	if err != nil{
		fmt.Println(err)
		return 0
	}
	rs.Close()
	return 1
}


func GetFromRedis(key string) string{
	rs, connectErr = connectRedis()

	value, err := redis.String(rs.Do("GET", key))

	if err != nil{
		fmt.Println(err)
	}
	rs.Close()
	return value
}


















