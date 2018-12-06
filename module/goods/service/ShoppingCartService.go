package service

import (
	"sync"
	"github.com/gin-gonic/gin"
	"shop-web/module/goods/dao"
	"shop-web/common/commonUtils"
	"shop-web/module/goods/protocol"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"shop-web/module/goods/model"
)

var ShoppingCartService = shoppingCartService{
	Mutex: &sync.Mutex{},
}

type shoppingCartService struct {
	Mutex *sync.Mutex
}

func (shoppingCartService *shoppingCartService)AddCart(c *gin.Context)  {

	shoppingCartService.Mutex.Lock()
	defer shoppingCartService.Mutex.Unlock()

	var err error

	userId := c.GetString("userId")
	if userId == "" {
		commonUtils.CreateErrorParams(c)
		return
	}

	var shoppingCartReq = protocol.ShoppingCartReq{}
	if err = c.ShouldBindWith(&shoppingCartReq,  binding.JSON); err != nil {
		commonUtils.CreateErrorParams(c)
		Logger.Error(err)
		return
	}
	Logger.Info(shoppingCartReq)

	exitCart := dao.HasSameSku(shoppingCartReq.SkuId, userId)
	if exitCart.Id != 0 {
		err := dao.UpdateSkuCount(exitCart.Id, shoppingCartReq.GoodsCount)
		if err != nil {
			commonUtils.CreateError(c)
			Logger.Error(err)
			return
		}

	}else {
		cart := model.ShoppingCart{}
		cart.ProductId = shoppingCartReq.ProductId
		cart.CheckStatus = 1
		cart.GoodsCount = shoppingCartReq.GoodsCount
		cart.UserId = userId
		cart.Status = 1
		goods, err := dao.FindGoodsByGoodsId(shoppingCartReq.ProductId)
		if err != nil || goods.ProductId == 0 {
			commonUtils.CreateErrorWithMsg(c, "sku not exits")
			Logger.Error(err)
			return
		}

		err = dao.AddCart(cart)
		if err != nil {
			commonUtils.CreateError(c)
			Logger.Error(err)
			return
		}
	}

	commonUtils.CreateSuccess(c, nil)
}

func (*shoppingCartService) FindCartList(c *gin.Context) {

	var language = c.GetString("language")
	var userId = c.GetString("userId")

	cartList,err := dao.FindCartList(userId)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return 
	}

	if len(cartList) == 0 {
		commonUtils.CreateNotContent(c)
		return
	}

	shoppingCartRespList := make([]protocol.ShoppingCartResp, 0)
	for _,cart := range cartList {
		shoppingCartResp := protocol.ShoppingCartResp{}
		shoppingCartResp.Id = cart.Id
		shoppingCartResp.ProductId = cart.ProductId
		shoppingCartResp.SkuValues = cart.SkuValues
		shoppingCartResp.GoodsCount = cart.GoodsCount

		goods,_ := dao.FindGoodsByGoodsId(cart.ProductId)

		if language == "EN" {
			shoppingCartResp.Title = goods.EnTitle
			shoppingCartResp.SellPoint = goods.EnSellPoint

		}else {
			shoppingCartResp.Title = goods.Title
			shoppingCartResp.SellPoint = goods.SellPoint
		}

		shoppingCartResp.Price = goods.Price
		shoppingCartResp.MemberPrice = goods.MemberPrice
		shoppingCartResp.StockNumber = goods.Stock

		if cart.CheckStatus == 1  {
			shoppingCartResp.CheckStatus = true
		}else {
			shoppingCartResp.CheckStatus = false
		}

		shoppingCartResp.GoodsImages = goods.GoodsImages
		shoppingCartResp.Status = cart.Status

		shoppingCartRespList = append(shoppingCartRespList, shoppingCartResp)
	}

	Logger.Info(shoppingCartRespList)

	commonUtils.CreateSuccessByList(c, len(shoppingCartRespList), shoppingCartRespList)
}

func (shoppingCartService *shoppingCartService)DelCart(c *gin.Context)  {

	shoppingCartService.Mutex.Lock()
	defer shoppingCartService.Mutex.Unlock()

	cartId := c.Param("cartId")
	userId := c.GetString("userId")
	if userId == "" {
		commonUtils.CreateErrorParams(c)
		return
	}

	id,_ := strconv.ParseInt(cartId, 10 , 64)
	err := dao.DelCart(id,userId)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateErrorParams(c)
		return
	}
	commonUtils.CreateSuccess(c, nil)

}