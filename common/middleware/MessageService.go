package middleware

import (
	ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
	"sync"
	"shop-web/common/redis"
	"shop-web/common/log"
	"os"
	"strconv"
	"math/rand"
	"shop-web/module/user/utils"
	"github.com/gin-gonic/gin"
	"shop-web/common/commonUtils"
	"time"
)

const (
	appkey = "8af596c745e7ce9eedd7bcea54fb114c"
	content = "【TOH】Your verification code is "
)

var logger = log.NewLogger(os.Stdout)

var MessageService = &messageService{
	mutex: &sync.Mutex{},
}

type messageService struct {
	 mutex *sync.Mutex
}

func (messageService *messageService)SendCode(c *gin.Context)  {
	// 发送短信
	messageService.mutex.Lock()
	defer messageService.mutex.Unlock()

	rand.Seed(time.Now().Unix())
	phone := c.Param("phone")
	code := strconv.Itoa(rand.Intn(1000000))
	key := utils.UserCodeKey +  phone
	err := redis.RedisClient.Set(key, []byte(code),time.Minute*5).Err()
	if err != nil {
		logger.Infof("redis set code error", err.Error())
		return
	}
	time.Now()
	clnt := ypclnt.New(appkey)
	param := ypclnt.NewParam(2)
	logger.Info(phone)
	param[ypclnt.MOBILE] = phone
	param[ypclnt.TEXT] = content + code
	r := clnt.Sms().SingleSend(param)
	logger.Info(r)
	if r.Code == 0 {
		commonUtils.CreateSuccess(c, code)
		return
	}

	commonUtils.CreateError(c)

	//账户:clnt.User() 签名:clnt.Sign() 模版:clnt.Tpl() 短信:clnt.Sms() 语音:clnt.Voice() 流量:clnt.Flow()
}