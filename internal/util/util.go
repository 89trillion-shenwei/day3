package util

import model2 "day3/internal/model"

func Init() {
	//储存礼品码信息
	model2.RedisPool = model2.NewRedisPool(model2.RedisURL, 1)
	//储存领取信息
	model2.RedisPool1 = model2.NewRedisPool(model2.RedisURL, 2)
}
