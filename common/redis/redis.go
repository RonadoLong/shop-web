package redis

import (
	"github.com/go-redis/redis"
	"time"
	"github.com/sirupsen/logrus"
)

//声明一些全局变量
var (
	RedisClient *redis.Client
)

func InitRedis(host string, password string)  {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,  // use default DB
		IdleTimeout: 240 * time.Second,
	})


	log := logrus.New()
	pong, err := RedisClient.Ping().Result()
	if err != nil {
		log.Error("connec redis err ", err)
	}
	log.Info(pong)

	go PSubscribeEventHandel()
}
