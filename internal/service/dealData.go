package service

import (
	"day3/app/model"
	"day3/internal"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
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

// Message redis存储的信息
type Message struct {
	Description    string    //礼品描述
	List           []List    //礼品内容列表（物品，数量）
	AvailableTimes string    //可领取次数
	ValidPeriod    string    //有效期
	GiftCode       string    //礼品码
	ReceivedTimes  string    //已领取次数
	Creator        string    //创建者账号
	CreatTime      string    //创建时间
	GetList        []GetList //领取列表
}

// Creator 管理员信息
type Creator struct {
	CreaName string //管理员账号
}

// User 用户信息
type User struct {
	UserName string //用户账号
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

//结构体转json
func struct2json(me interface{}) []byte {
	byts, _ := json.Marshal(me)
	return byts
}

//json转结构体
func json2struct(byts []byte, message *Message) bool {
	err := json.Unmarshal(byts, message)
	if err != nil {
		return false
	}
	return true

}

// StrSet 创建数据
func (creator *Creator) StrSet(key string, message Message) error {
	c := model.NewRedisPool(model.RedisURL, 1).Get()
	defer c.Close()
	//将结构体转为json字符串在存入redis
	_, err2 := c.Do("SET", key, struct2json(message))
	if err2 != nil {
		return internal.InternalServiceError(err2.Error())
	} else {
		return nil
	}
}

// StrGet 查询数据，返回所有数据
func (creator *Creator) StrGet(key string) string {
	c := model.NewRedisPool(model.RedisURL, 1).Get()
	defer c.Close()
	res, _ := redis.String(c.Do("GET", key))
	return res
}

// StrUpdate 用户领取礼品时更新数据库，增加领取人列表，修改可领取次数和已领取次数，返回礼品列表
func (User *User) StrUpdate(key string) (string, error) {
	c := model.NewRedisPool(model.RedisURL, 1).Get()
	defer c.Close()
	//查询数据
	res, err := redis.String(c.Do("Get", key))
	if err != nil {
		return "查询礼品时出错", internal.InternalServiceError(err.Error())
	} else {
		var byts []byte
		byts = []byte(res)
		message := Message{}
		//json转结构体成功
		if json2struct(byts, &message) {
			fmt.Println("success")
		}
		getList := new(GetList)
		//用户名
		getList.GetorName = User.UserName
		//领取时间
		getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
		//判断领取时间是否超出有效期
		if String2Time(getList.GetTime) >= String2Time(message.ValidPeriod) {
			return "该礼品码已过期", internal.KeyExpiredError("该礼品码已过期")
		}
		//判断该用户是否已经使用过该礼品码
		if findUser(message.GetList, User.UserName) {
			return "你已使用过该礼品码", internal.UserHasEeceivedError("你已使用过该礼品码")
		}
		message.GetList = append(message.GetList, *getList)
		//可领取次数
		av, _ := strconv.Atoi(message.AvailableTimes)
		if av == 0 {
			return "该礼品码已被领取完毕", internal.NoGiftError("该礼品码已被领取完毕")
		}
		av -= 1
		//已领取次数
		re, _ := strconv.Atoi(message.ReceivedTimes)
		re += 1
		message.ReceivedTimes = strconv.Itoa(re)
		message.AvailableTimes = strconv.Itoa(av)
		//提交更改后的数据
		_, err := c.Do("SET", key, struct2json(message))
		if err != nil {
			return "", internal.InternalServiceError(err.Error())
		} else {
			fmt.Println("set ok.")
		}
		//返回礼品内容
		return string("您将获得的礼品有：" + string(struct2json(message.List))), nil
	}
}

// CheckKey 判断数据是否存在
func CheckKey(key string) bool {
	c := model.NewRedisPool(model.RedisURL, 1).Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false
	} else {
		return exist
	}
}

//判断用户是否已经领取过礼品
func findUser(list []GetList, name string) bool {
	for _, item := range list {
		if item.GetorName == name {
			return true
		}
	}
	return false
}
