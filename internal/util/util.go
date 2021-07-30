package util

import (
	"day3/internal/model"
)

func Init() {
	//储存礼品码信息
	model.RedisPool = model.NewRedisPool(model.RedisURL, 1)
	//储存领取信息
	model.RedisPool1 = model.NewRedisPool(model.RedisURL, 2)
}
