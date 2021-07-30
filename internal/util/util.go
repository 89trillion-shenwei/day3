package util

import (
	"day3/internal/service"
)

func Init() {
	//储存礼品码信息
	service.RedisPool = service.NewRedisPool(service.RedisURL, 1)
	//储存领取信息
	service.RedisPool1 = service.NewRedisPool(service.RedisURL, 2)
}
