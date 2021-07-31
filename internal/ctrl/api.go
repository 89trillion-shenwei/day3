package ctrl

import (
	"day3/internal"
	"day3/internal/handler"
	"day3/internal/model"
	"day3/internal/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SetStrApi 录入礼品
func SetStrApi(c *gin.Context) (string, error) {
	message := new(model.Message)
	mess := new(model.Mess)
	//礼品描述
	message.Description = c.PostForm("Description")
	if message.Description == "" {
		return "", internal.IsEmptyError("礼品描述不能为空")
	}
	//礼品码类型
	message.CodeType = c.PostForm("CodeType")
	if message.CodeType != "1" && message.CodeType != "2" && message.CodeType != "3" {
		return "", internal.CodeTypeError("礼品码类型不符合规范")
	}
	//可领取用户
	message.CanGetUser = c.PostForm("CanGetUser")
	//创建者账号
	message.Creator = c.PostForm("Creator")
	if message.Creator == "" {
		return "", internal.IsEmptyError("创建者账号不能为空")
	}
	//创建时间,默认为当前时间
	message.CreatTime = time.Now().Format("2006-01-02 15:04:05")
	//有效期
	message.ValidPeriod = c.PostForm("ValidPeriod")
	if message.ValidPeriod == "" {
		return "", internal.IsEmptyError("有效期不能为空")
	}
	if model.String2Time(message.CreatTime) >= model.String2Time(message.ValidPeriod) {
		return "", internal.ValidPeriodError("有效期不能小于当前时间")
	}
	//可领取次数
	mess.AvailableTimes = c.PostForm("AvailableTimes")
	if mess.AvailableTimes == "" {
		return "", internal.IsEmptyError("可领取次数不能为空")
	}
	//已领取次数，默认为0
	mess.ReceivedTimes = "0"
	//礼品列表，用","分割
	listStr := c.PostForm("List")
	if listStr == "" {
		return "", internal.IsEmptyError("礼品列表不能为空")
	}
	s := strings.Split(listStr, ",")
	for i := 0; i < len(s)/2; i++ {
		list := new(model.List)
		list.Name = s[i*2]
		list.Amount = s[i*2+1]
		message.List = append(message.List, *list)
	}
	creator := new(service.Creator)
	creator.CreaName = c.PostForm("Creator")
	re := handler.Set(*message, *mess, *creator)
	return re, nil
}

// GetStrApi 根据礼品码查询礼品
func GetStrApi(c *gin.Context) (s1, s2 string, err error) {
	//礼品码
	key := c.PostForm("key")
	if key == "" {
		return "", "", internal.IsEmptyError("礼品码不能为空")
	}
	if len(key) != 8 {
		return "", "", internal.LenFalseError("礼品码不合法")
	}
	creator := new(service.Creator)
	re1, re2, err := handler.Get(key, *creator)
	if err != nil {
		return "", "", err
	}
	return re1, re2, nil
}

func StrUpdateApi(c *gin.Context) (string, error) {
	user := new(service.User)
	key := c.PostForm("key")
	if key == "" {
		return "", internal.IsEmptyError("礼品码不能为空")
	}
	if len(key) != 8 {
		return "", internal.LenFalseError("礼品码不合法")
	}
	user.UserName = c.PostForm("username")
	if user.UserName == "" {
		return "", internal.IsEmptyError("用户名不能为空")
	}
	re, err := handler.Update(*user, key)
	if err != nil {
		return "", err
	}
	return re, nil
}

type Api func(c *gin.Context) (string, error)
type Api1 func(c *gin.Context) (string, string, error)

func ReturnData(api Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := api(c)
		if err != nil {
			globalError := err.(internal.GlobalError)
			c.JSON(globalError.Status, globalError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"查询结果": data,
		})
	}
}

func ReturnData1(api Api1) gin.HandlerFunc {
	return func(c *gin.Context) {
		data1, data2, err := api(c)
		if err != nil {
			globalError := err.(internal.GlobalError)
			c.JSON(globalError.Status, globalError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"礼品信息": data1,
			"领取信息": data2,
		})
	}
}
