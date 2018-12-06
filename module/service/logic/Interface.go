package logic

import (
	"github.com/gin-gonic/gin"
	"os"
	"shop-web/common/log"
)

type ServiceInterface interface {
	SetContext()
	SaveService(c *gin.Context)
	FindService(c *gin.Context)
}


var (
	Cname   = "service"
	Logger  = log.NewLogger(os.Stdout)
)