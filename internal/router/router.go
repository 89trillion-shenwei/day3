package router

import (
	"day3/internal/ctrl"
	"github.com/gin-gonic/gin"
)

func SetStrRouter() *gin.Engine {
	router := gin.Default()
	//录入数据
	router.POST("/SetStr", ctrl.SetStrApi)
	//查询数据
	router.POST("/GetStr", ctrl.GetStrApi)
	//领取礼品
	router.POST("/UpdateStr", ctrl.StrUpdate)
	return router
}
