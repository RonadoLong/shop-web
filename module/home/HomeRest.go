package home

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/home/homeService"
)

var HomeRest = homeRest{}

type homeRest struct {
	
}

func (*homeRest)FindHomeNav(c *gin.Context) {
	home.HomeService.FindHomeNav(c)
}


func (*homeRest)FindHomeList(c *gin.Context) {
	home.HomeService.FindHomeList(c)
}