package handler

import (
	"day3/internal/service"
	"testing"
)

//测试创建数据方法，若redis数据库里可以找到数据则success
func TestSet(t *testing.T) {
	message := service.Message{}
	creator := service.Creator{}
	message.Description = "测试一"
	creator.CreaName = "测试员1"
	key := Set(message, creator)
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
	key := "a8e46a16"
	creator := service.Creator{}
	creator.CreaName = "测试员1"
	got, _ := Get(key, creator)
	if got != "" {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}

//测试领取礼品方法
func TestUpdate(t *testing.T) {
	key := "a8e46a16"
	user := service.User{}
	user.UserName = "用户一"
	got, _ := Update(user, key)
	want := "该礼品码已被领取完毕"
	if got == want {
		t.Log("success")
		return
	} else {
		t.Log("failed")
		t.Error("failed")
	}
}
