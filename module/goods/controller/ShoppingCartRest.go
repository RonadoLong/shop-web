package controller

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/goods/service"
)

var ShoppingCartRest = shoppingCartRest{}

type shoppingCartRest struct {

}

func (*shoppingCartRest) AddCart(c *gin.Context)  {
	service.ShoppingCartService.AddCart(c)
}

func (*shoppingCartRest) FindCartList(c *gin.Context) {
	service.ShoppingCartService.FindCartList(c)
}

func (*shoppingCartRest) DelCart(c *gin.Context)  {
	service.ShoppingCartService.DelCart(c)
}