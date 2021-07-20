package ctrl

import (
	"day3/internal"
	"day3/internal/handler"
	"day3/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// SetStrApi 录入礼品
func SetStrApi(c *gin.Context) error {
	message := new(service.Message)
	//礼品描述
	message.Description = c.PostForm("Description")
	if message.Description == "" {
		return internal.IsEmptyError("礼品描述不能为空")
	}
	//创建者账号
	message.Creator = c.PostForm("Creator")
	if message.Creator == "" {
		return internal.IsEmptyError("创建者账号不能为空")
	}
	//创建时间,默认为当前时间
	message.CreatTime = time.Now().Format("2006-01-02 15:04:05")
	//可领取次数
	message.AvailableTimes = c.PostForm("AvailableTimes")
	if message.AvailableTimes == "" {
		return internal.IsEmptyError("可领取次数不能为空")
	}
	//已领取次数，默认为0
	message.ReceivedTimes = "0"
	//礼品列表，用","分割
	listStr := c.PostForm("List")
	if listStr == "" {
		return internal.IsEmptyError("礼品列表不能为空")
	}
	s := strings.Split(listStr, ",")
	for i := 0; i < len(s)/2; i++ {
		list := new(service.List)
		list.Name = s[i*2]
		list.Amount = s[i*2+1]
		message.List = append(message.List, *list)
	}
	//有效期
	message.ValidPeriod = c.PostForm("ValidPeriod")
	if message.ValidPeriod == "" {
		return internal.IsEmptyError("有效期不能为空")
	}
	if service.String2Time(message.CreatTime) >= service.String2Time(message.ValidPeriod) {
		return internal.ValidPeriodError("有效期不能小于当前时间")
	}
	creator := new(service.Creator)
	creator.CreaName = c.PostForm("Creator")
	re := handler.Set(*message, *creator)
	//返回礼品码
	c.String(http.StatusOK, re)
	return nil
}

// GetStrApi 根据礼品码查询礼品
func GetStrApi(c *gin.Context) error {
	//礼品码
	key := c.PostForm("key")
	if key == "" {
		return internal.IsEmptyError("礼品码不能为空")
	}
	if len(key) != 8 {
		return internal.LenFalseError("礼品码不合法")
	}
	creator := new(service.Creator)
	re, _ := handler.Get(key, *creator)
	c.String(http.StatusOK, re)
	return nil
}

func StrUpdate(c *gin.Context) error {
	user := new(service.User)
	key := c.PostForm("key")
	if key == "" {
		return internal.IsEmptyError("礼品码不能为空")
	}
	if len(key) != 8 {
		return internal.LenFalseError("礼品码不合法")
	}
	user.UserName = c.PostForm("username")
	if user.UserName == "" {
		return internal.IsEmptyError("用户名不能为空")
	}
	re, _ := handler.Update(*user, key)
	c.String(http.StatusOK, re)
	return nil
}
