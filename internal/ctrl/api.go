package ctrl

import (
	"day3/app/model"
	"day3/internal/handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// SetStrApi 录入礼品
func SetStrApi(c *gin.Context) {
	message := new(model.Message)
	//礼品描述
	message.Description = c.PostForm("Description")
	//创建者账号
	message.Creator = c.PostForm("Creator")
	//创建时间,默认为当前时间
	message.CreatTime = time.Now().Format("2006-01-02 15:04:05")
	//可领取次数
	message.AvailableTimes = c.PostForm("AvailableTimes")
	//已领取次数，默认为0
	message.ReceivedTimes = "0"
	//礼品列表，用","分割
	listStr := c.PostForm("List")
	s := strings.Split(listStr, ",")
	for i := 0; i < len(s)/2; i++ {
		list := new(model.List)
		list.Name = s[i*2]
		list.Amount = s[i*2+1]
		message.List = append(message.List, *list)
	}
	//有效期
	message.ValidPeriod = c.PostForm("ValidPeriod")
	creator := new(model.Creator)
	creator.CreaName = c.PostForm("Creator")
	re := handler.Set(*message, *creator)
	//返回礼品码
	c.String(http.StatusOK, re)
}

// GetStrApi 根据礼品码查询礼品
func GetStrApi(c *gin.Context) {
	//礼品码
	key := c.PostForm("key")
	creator := new(model.Creator)
	re, _ := handler.Get(key, *creator)
	c.String(http.StatusOK, re)
}

func StrUpdate(c *gin.Context) {
	user := new(model.User)
	key := c.PostForm("key")
	user.UserName = c.PostForm("username")
	re, _ := handler.Update(*user, key)
	c.String(http.StatusOK, re)
}
