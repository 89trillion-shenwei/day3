package main

import (
	"day3/app/http"
	"day3/internal/util"
)

func main() {
	//初始化redis池
	util.Init()
	//启动服务
	http.Start()
}
