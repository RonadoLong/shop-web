package engine

import (
	"os"
	"shop-web/common/jwt"
	"shop-web/common/log"
	"shop-web/common/middleware"
	"shop-web/module/user/controller"
	"shop-web/module/user/model"
	"time"

	"github.com/gin-gonic/gin"

	Goods "shop-web/module/goods/controller"
	"shop-web/module/home"
	News "shop-web/module/news/controller"
	Order "shop-web/module/order/controller"
	Service "shop-web/module/service/controller"
)

var Logger = log.NewLogger(os.Stdout)

func ClientEngine() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	AuthMiddleware := &jwt.GinJWTMiddleware{
		PubKeyFile:  "resource/pub_key.pem",
		PrivKeyFile: "resource/pri_key.pem",
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour * 24 * 3000,
		MaxRefresh:  time.Hour,
		Authenticator: func(loginParams interface{}, c *gin.Context) (*model.User, bool) {
			return controller.UserRest.Login(loginParams, c)
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
	api := router.Group("/api/client")

	/**  global interceptor */
	api.Use(AuthMiddleware.MiddlewareParseUser)
	api.Use(middleware.RequestInterceptor)

	/**  user logic */
	userGroup := api.Group("/user")
	userGroup.POST("/login", AuthMiddleware.LoginHandler)
	userGroup.GET("/getCode/:phone", controller.UserRest.SendCode)

	auth := userGroup.Group("/auth")
	auth.Use(AuthMiddleware.MiddlewareFunc())
	{
		auth.POST("/update", controller.UserRest.UpdateUserInfo)
		auth.GET("/userInfo", controller.UserRest.GetUserInfo)
		auth.GET("/saveAgreement", controller.UserRest.SaveAgreement)
		auth.GET("/getUserAgreement", controller.UserRest.GetUserAgreement)
		auth.GET("/refresh_token", AuthMiddleware.RefreshHandler)
		auth.GET("/getUserAddress", controller.UserRest.GetUserAddress)
		auth.POST("/saveUserAddress", controller.UserRest.SaveUserAddress)
		auth.GET("/getUserIntegralFlow/:pageNum/:pageSize", controller.UserRest.GetUserIntegralFlow)
	}

	/** home logic */
	homeGroup := api.Group("/home")
	homeGroup.GET("/headers", home.HomeRest.FindHomeNav)
	homeGroup.GET("/list", home.HomeRest.FindHomeList)

	/** news logic  */
	news := api.Group("/news")
	news.GET("/nav/list", News.NewsRest.FindNewSNavList)
	news.GET("/list/:category/:pageNum/:pageSize", News.NewsRest.FindNewsList)
	news.GET("/homeList/:pageNum/:pageSize", News.NewsRest.FindNewsAllList)
	news.GET("/detail/:newsId", News.NewsRest.FindNewsDetail)
	news.GET("/like/:newsId", News.NewsRest.AddLikeById)

	/** video logic  */
	video := api.Group("/video")
	video.GET("/list/:category/:pageNum/:pageSize", News.VideoRest.FindVideoList)
	video.GET("/homeList/:pageNum/:pageSize", News.VideoRest.FindVideoAllList)
	video.GET("/detail/:videoId", News.VideoRest.FindVideoDetail)
	video.GET("categoryList", News.VideoRest.FindCategoryList)

	/** goods logic  */
	goods := api.Group("/goods")
	goods.GET("/nav/list", Goods.GoodsRest.FindGoodsNavList)
	goods.GET("/list/:category/:pageNum/:pageSize", Goods.GoodsRest.FindGoodsList)
	goods.GET("/homeList/:pageNum/:pageSize", Goods.GoodsRest.FindGoodsAllList)
	goods.GET("/detail/:goodsId", Goods.GoodsRest.FindGoodsDetail)

	/** cart logic */
	cart := api.Group("/cart")
	cart.Use(AuthMiddleware.MiddlewareFunc())
	cart.POST("/confirm", Goods.ShoppingCartRest.AddCart)
	cart.GET("/list", Goods.ShoppingCartRest.FindCartList)
	cart.DELETE("/del/:cartId", Goods.ShoppingCartRest.DelCart)

	/** order logic */
	order := api.Group("/order")
	order.Use(AuthMiddleware.MiddlewareFunc())
	order.POST("/create", Order.OrderInfoRest.CreateOrder)
	order.POST("/list", Order.OrderInfoRest.FindOrderInfoList)
	order.PUT("/cancel/:orderId", Order.OrderInfoRest.CancelOrder)

	other := api.Group("/other")
	other.POST("/callback", Order.OrderInfoRest.ReceiveIPN)

	/** logic logic */
	service := api.Group("/service")
	service.GET("/category/list", Service.ServiceRest.FindServiceCategoryList)
	service.GET("/findServicePaymentList", Service.ServiceRest.FindServicePaymentList)
	service.GET("/area/list/:areaname", Service.ServiceRest.FindAreaListByParentId)
	service.POST("/room/list/:pageNum/:pageSize", Service.ServiceRest.FindServiceRoomByCategory)

	serviceAuth := service.Use(AuthMiddleware.MiddlewareFunc())
	serviceAuth.POST("/auth/save", Service.ServiceRest.AddServiceRoom)
	serviceAuth.POST("/job/auth/save", Service.ServiceRest.SaveServiceJob)
	serviceAuth.GET("/auth/self/:status/:pageNum/:pageSize", Service.ServiceRest.FindSelfService)

	return router
}
