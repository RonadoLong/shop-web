package middleware

import (
	"github.com/gin-gonic/gin"
	"shop-web/common/redis"
	"shop-web/common/commonUtils"
	"time"
	"shop-web/conf"
)

const (
	RATELIMITER_KEY = "RATELIMITER:"
)

func RequestInterceptor(c *gin.Context) {

	user := "guest"
	if c.GetString("userId") != "" {
		user = c.GetString("userId")
	}

	logger.Info("request user ===== ", user)
	logger.Info("request user ip ===== ", c.ClientIP())
	logger.Info("request user path ===== ", c.HandlerName())
	logger.Info("request user params ===== ", c.Params)
	logger.Info("request user deviceId ===== ", c.GetHeader("deviceId"))
	logger.Info("request user phoneModel ===== ", c.GetHeader("phoneModel"))
	logger.Info("request user language ===== ", c.GetHeader("language"))
	logger.Info("request user version ===== ", c.GetHeader("version"))

	if c.GetHeader("language") == "" {
		c.Set("language", "CN")
	}else {
		c.Set("language", c.GetHeader("language"));
	}

	if conf.GetConfig().GetString("config") == "dev" {
		rateLimier(c)
	}
}

func rateLimier(c *gin.Context) {
	key := RATELIMITER_KEY + c.ClientIP()
	ret, err := redis.RedisClient.Exists(key).Result()
	if err != nil {
		logger.Error(err)
	}

	if int(ret) > 0 {
		count, _ := redis.RedisClient.Get(key).Int64()
		if count >= 5 {
			commonUtils.CreateErrorRequest(c)
			return
		}
		redis.RedisClient.Incr(key)
	} else {
		redis.RedisClient.Set(key, int64(1), time.Second)
	}

}
