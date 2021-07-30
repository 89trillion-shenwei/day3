package model

import (
	"day3/internal"
	"day3/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	wg       sync.WaitGroup
	lockChan = make(chan struct{}, 1)
)

// 如果lockChan中为空则阻塞
func getLock() {
	<-lockChan
}

// 重新填充lockChan
func releaseLock() {
	lockChan <- struct{}{}
}

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

// Creator 管理员信息
type Creator struct {
	CreaName string //管理员账号
}

// User 用户信息
type User struct {
	UserName string //用户账号
}

// StrSet 创建数据
func (creator *Creator) StrSet(key string, message Message, mess Mess) error {
	c1 := RedisPool.Get()
	c2 := RedisPool1.Get()
	defer c1.Close()
	defer c2.Close()
	time := service.String2Time(message.ValidPeriod) - service.String2Time(message.CreatTime)
	//将结构体转为json字符串在存入redis
	_, err2 := c1.Do("SET", key, service.Struct2json(message))
	//将结构体转为json字符串在存入redis
	_, err3 := c2.Do("SET", key, service.Struct2json(mess))
	//根据时间戳设置过期时间
	n1, _ := c1.Do("EXPIRE", key, time)
	n2, _ := c2.Do("EXPIRE", key, time)
	if n1 == int64(1) && n2 == int64(1) {
		fmt.Println("success")
	}
	if err2 != nil {
		return internal.InternalServiceError(err2.Error())
	} else {
		return nil
	}
	if err3 != nil {
		return internal.InternalServiceError(err3.Error())
	} else {
		return nil
	}
}

// GetGiftCodeInformation 获取礼品码信息
func (creator *Creator) GetGiftCodeInformation(key string) string {
	c := RedisPool.Get()
	defer c.Close()
	res, _ := redis.String(c.Do("GET", key))
	return res
}

// GetGiftCollectionInformation 获取领取信息
func (creator *Creator) GetGiftCollectionInformation(key string) string {
	c := RedisPool1.Get()
	defer c.Close()
	res, _ := redis.String(c.Do("GET", key))
	return res
}

// StrUpdate 用户领取礼品时更新数据库，增加领取人列表，修改可领取次数和已领取次数，返回礼品列表
func (User *User) StrUpdate(key string) (string, error) {
	c1 := RedisPool.Get()
	c2 := RedisPool1.Get()
	defer c1.Close()
	defer c2.Close()
	//上锁
	releaseLock()
	//查询数据
	res1, err1 := redis.String(c1.Do("Get", key))
	res2, err2 := redis.String(c2.Do("Get", key))
	if err1 != nil || err2 != nil {
		//解锁
		getLock()
		return "查询礼品时出错", internal.InternalServiceError(err1.Error() + err2.Error())
	} else {
		var byts1, byts2 []byte
		byts1 = []byte(res1) //礼品码信息
		byts2 = []byte(res2) //领取信息
		mess := Mess{}
		message := Message{}
		//json转结构体成功
		if Json2struct(byts1, &message) && Json2struct1(byts2, &mess) {
			fmt.Println("success")
		}
		//如果礼品可领取次数为0，退出
		if mess.AvailableTimes == "0" {
			//解锁
			getLock()
			return "", internal.NoGiftError("礼品已领完")
		}
		//指定用户一次性消耗
		if message.CodeType == "1" {
			//如果用户是指定用户
			if message.CanGetUser == User.UserName {
				getList := new(GetList)
				//用户名
				getList.GetorName = User.UserName
				//领取时间
				getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
				mess.GetList = append(mess.GetList, *getList)
				//可领取次数
				av, _ := strconv.Atoi(mess.AvailableTimes)
				av -= 1
				//已领取次数
				re, _ := strconv.Atoi(mess.ReceivedTimes)
				re += 1
				mess.ReceivedTimes = strconv.Itoa(re)
				mess.AvailableTimes = strconv.Itoa(av)
				//提交更改后的数据
				_, err := c2.Do("SET", key, service.Struct2json(mess))
				if err != nil {
					//解锁
					getLock()
					return "", internal.InternalServiceError(err.Error())
				} else {
					fmt.Println("set ok.")
				}
				//incr key
				_, err1 := redis.Int64(c2.Do("INCR", "key"))
				if err1 != nil {
					log.Println("INCR failed:", err)
					//解锁
					getLock()
					return "", internal.InternalServiceError(err1.Error())
				}
			} else {
				//解锁
				getLock()
				return "", internal.NoCanGetUserError("非指定用户")
			}
		} else if message.CodeType == "2" { //不指定用户,限制次数
			getList := new(GetList)
			//用户名
			getList.GetorName = User.UserName
			//领取时间
			getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
			//判断该用户是否已经使用过该礼品码
			if FindUser(mess.GetList, User.UserName) {
				//解锁
				getLock()
				return "", internal.UserHasEeceivedError("你已使用过该礼品码")
			}
			mess.GetList = append(mess.GetList, *getList)
			//可领取次数
			av, _ := strconv.Atoi(mess.AvailableTimes)
			av -= 1
			//已领取次数
			re, _ := strconv.Atoi(mess.ReceivedTimes)
			re += 1
			mess.ReceivedTimes = strconv.Itoa(re)
			mess.AvailableTimes = strconv.Itoa(av)
			//提交更改后的数据
			_, err := c2.Do("SET", key, service.Struct2json(mess))
			if err != nil {
				//解锁
				getLock()
				return "", internal.InternalServiceError(err.Error())
			} else {
				fmt.Println("set ok.")
			}
			//incr key
			_, err1 := redis.Int64(c2.Do("INCR", "key"))
			if err1 != nil {
				log.Println("INCR failed:", err)
				//解锁
				getLock()
				return "", internal.InternalServiceError(err1.Error())
			}
		} else { //不指定用户,不限制次数
			getList := new(GetList)
			//用户名
			getList.GetorName = User.UserName
			//领取时间
			getList.GetTime = time.Now().Format("2006-01-02 15:04:05")
			//判断该用户是否已经使用过该礼品码
			if FindUser(mess.GetList, User.UserName) {
				//解锁
				getLock()
				return "你已使用过该礼品码", internal.UserHasEeceivedError("你已使用过该礼品码")
			}
			mess.GetList = append(mess.GetList, *getList)
			//已领取次数
			re, _ := strconv.Atoi(mess.ReceivedTimes)
			re += 1
			mess.ReceivedTimes = strconv.Itoa(re)
			mess.AvailableTimes = "99999999"
			//提交更改后的数据
			_, err := c2.Do("SET", key, service.Struct2json(mess))
			if err != nil {
				//解锁
				getLock()
				return "", internal.InternalServiceError(err.Error())
			} else {
				fmt.Println("set ok.")
			}
			//incr key
			_, err1 := redis.Int64(c2.Do("INCR", "key"))
			if err1 != nil {
				log.Println("INCR failed:", err)
				//解锁
				getLock()
				return "", internal.InternalServiceError(err1.Error())
			}
		}
		//解锁
		getLock()
		//返回礼品内容
		return string("您将获得的礼品有：" + string(service.Struct2json(message.List))), nil
	}
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

// FindUser 判断用户是否已经领取过礼品
func FindUser(list []GetList, name string) bool {
	for _, item := range list {
		if item.GetorName == name {
			return true
		}
	}
	return false
}
