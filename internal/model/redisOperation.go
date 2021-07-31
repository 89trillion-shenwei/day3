package model

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

// List 物品列表
type List struct {
	Name   string //物品名
	Amount string // 物品数量
}

type GetList struct {
	GetorName string //领取人用户名
	GetTime   string //领取时间
}

// Message redis存储的礼品码信息
type Message struct {
	Description string //礼品描述
	CodeType    string //礼品码类型
	List        []List //礼品内容列表（物品，数量）
	ValidPeriod string //有效期
	GiftCode    string //礼品码
	CanGetUser  string //允许领取用户
	Creator     string //创建者账号
	CreatTime   string //创建时间
}

// Mess redis存储的领取信息
type Mess struct {
	AvailableTimes string    //可领取次数
	ReceivedTimes  string    //已领取次数
	GiftCode       string    //礼品码
	key            string    //计数
	GetList        []GetList //领取列表
}

// CheckKey 判断数据是否存在
func CheckKey(key string) bool {
	c := RedisPool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false
	} else {
		return exist
	}
}

// FindUser 判断用户是否已经领取过礼品
func FindUser(list []GetList, name string) bool {
	for _, item := range list {
		if item.GetorName == name {
			return true
		}
	}
	return false
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

// Json2struct json转结构体
func Json2struct(byts []byte, message *Message) bool {
	err := json.Unmarshal(byts, message)
	if err != nil {
		return false
	}
	return true

}

// Json2struct1 json转结构体
func Json2struct1(byts []byte, message *Mess) bool {
	err := json.Unmarshal(byts, message)
	if err != nil {
		return false
	}
	return true
}
