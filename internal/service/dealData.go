package service

import (
	"day3/internal/model"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
)

// CheckKey 判断数据是否存在
func CheckKey(key string) bool {
	c := model.RedisPool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false
	} else {
		return exist
	}
}

//字符串转时间戳
func String2Time(s string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	if err != nil {
		return 0
	}
	unixTime := theTime.Unix()
	return unixTime
}

// Struct2json 结构体转json
func Struct2json(me interface{}) []byte {
	byts, _ := json.Marshal(me)
	return byts
}
