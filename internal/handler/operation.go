package handler

import (
	"crypto/md5"
	"day3/app/model"
	"day3/internal"
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
func Set(message model.Message, creator model.Creator) string {
	key := CreatePasswd()
	//如果物品码不在数据库中，正常情况下是不可能一样
	if !model.CheckKey(key) {
		message.GiftCode = key
		creator.StrSet(key, message)
		return key
	} else {
		//数据库已经有了就重新生成一个
		key1 := CreatePasswd()
		message.GiftCode = key1
		creator.StrSet(key1, message)
		return key1
	}
}

// Get 查询数据
func Get(key string, creator model.Creator) (string, error) {
	//如果key存在
	if model.CheckKey(key) {
		return creator.StrGet(key), nil
	} else {
		fmt.Println("无此数据")
		return "无此数据", internal.NoKeyError("礼品码不存在")
	}
}

// Update 用户领取礼品，更新数据
func Update(user model.User, key string) (string, error) {
	if model.CheckKey(key) {
		return user.StrUpdate(key)
	} else {
		return "礼品码不存在，请重新输入！", internal.NoKeyError("礼品码不存在")
	}

}
