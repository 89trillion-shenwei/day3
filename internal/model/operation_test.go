package model

import (
	"day3/internal/service"
	"day3/internal/util"
	"testing"
	"time"
)

//测试创建数据方法，若redis数据库里可以找到数据则success
func TestSet(t *testing.T) {
	util.Init()
	message := service.Message{}
	mess := service.Mess{}
	creator := service.Creator{}
	message.Description = "测试一"
	creator.CreaName = "测试员1"
	message.ValidPeriod = "2022-01-02 15:04:05"
	list1 := new(service.List)
	list2 := new(service.List)
	list1.Name = "1001"
	list1.Amount = "4"
	list2.Name = "1002"
	list1.Amount = "8"
	message.List = append(message.List, *list1)
	message.List = append(message.List, *list2)
	message.GiftCode = "3"
	message.CreatTime = time.Now().Format("2006-01-02 15:04:05")
	mess.AvailableTimes = "100"
	key := Set(message, mess, creator)
	if service.CheckKey(key) {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}

//测试查询数据方法
func TestGet(t *testing.T) {
	util.Init()
	key := "45e72e99"
	creator := service.Creator{}
	creator.CreaName = "测试员1"
	got1, got2, _ := Get(key, creator)
	if got1 != "" && got2 != "" {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}

//测试领取礼品方法
func TestUpdate(t *testing.T) {
	util.Init()
	key := "45e72e99"
	user := service.User{}
	user.UserName = "用户一"
	got, _ := Update(user, key)
	if got != "" {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}
