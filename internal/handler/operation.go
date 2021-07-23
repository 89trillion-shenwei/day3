package handler

import (
	"crypto/md5"
	"day3/internal"
	"day3/internal/service"
	"fmt"
	"io"
	"time"
)

type Creator struct {
	CreaName string //管理员账号
}

// CreatePasswd 随机生成物品码
func CreatePasswd() string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "znxgdukahedue")
	io.WriteString(h, t.String())
	passwd := fmt.Sprintf("%x", h.Sum(nil))
	return passwd[:8]
}

// Set 存数据
func Set(message service.Message, mess service.Mess, creator service.Creator) string {
	key := CreatePasswd()
	//如果物品码不在数据库中，正常情况下是不可能一样
	if !service.CheckKey(key) {
		message.GiftCode = key
		mess.GiftCode = key
		creator.StrSet(key, message, mess)
		return key
	} else {
		//数据库已经有了就重新生成一个
		key1 := CreatePasswd()
		message.GiftCode = key1
		mess.GiftCode = key1
		creator.StrSet(key1, message, mess)
		return key1
	}
}

// Get 查询数据
func Get(key string, creator service.Creator) (s1, s2 string, err error) {
	//如果key存在
	if service.CheckKey(key) {
		return creator.GetGiftCodeInformation(key), creator.GetGiftCollectionInformation(key), nil
	} else {
		return "", "", internal.NoKeyError("礼品码不存在")
	}
}

// Update 用户领取礼品，更新数据
func Update(user service.User, key string) (string, error) {
	if service.CheckKey(key) {
		return user.StrUpdate(key)
	} else {
		return "礼品码不存在，请重新输入！", internal.NoKeyError("礼品码不存在")
	}

}
