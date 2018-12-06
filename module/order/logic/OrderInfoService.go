package logic

import (
	"fmt"
	"os"
	"shop-web/common/commonUtils"
	"shop-web/common/dbutil"
	"shop-web/common/log"
	"shop-web/common/protocol"
	"shop-web/common/redis"
	"shop-web/module/order/dao"
	protocol2 "shop-web/module/order/protocol"
	"sort"
	"strings"
	"sync"

	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"shop-web/module/service/model"
	"shop-web/module/service/logic"
	model2 "shop-web/module/order/model"
	UserDao "shop-web/module/user/dao"
	model3 "shop-web/module/user/model"
	"time"
)

var logger = log.NewLogger(os.Stdout)

var OrderInfoService = orderInfoService{
	Mutex: &sync.Mutex{},
}


type message struct {
	OrderId string
	Note note
}

type note struct {
	ServiceId string
	//0 为商品 1 为服务
	Type string
	//支付方式
	PayType string
	SelectId int
}

var worker = make(chan *message, 1024)

var saveWorker = make(chan interface{}, 1024)

func (*orderInfoService)InitDetail()  {
	OrderInfoService.startSaveOrderflow()
}

func (js *orderInfoService) SetContext() {
	js.context = dbutil.NewDBContext()
	js.c = js.context.DBCollection("order_pay_flow")
}

type orderInfoService struct {
	Mutex   *sync.Mutex
	context *dbutil.MongoContext
	c       *mgo.Collection
}

var token = "0c54c847ce560c8cd64575db9f8a13d7ed752f2de5c0a080ce228e1610fd6d8d"

func (orderService *orderInfoService) CreateOrder(c *gin.Context) {

	orderService.Mutex.Lock()
	defer orderService.Mutex.Unlock()

	confirmOrderReq := protocol.ConfirmOrderReq{}
	if err := c.ShouldBindJSON(&confirmOrderReq); err != nil {
		commonUtils.CreateErrorParams(c)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		commonUtils.CreateErrorParams(c)
		return
	}

	confirmOrderReq.UserId = userId
	//language := c.GetString("language")

	/** 检查参数 */
	if confirmOrderReq.AddressId == "" || len(confirmOrderReq.OrderGoodsReqList) == 0 {
		commonUtils.CreateErrorParams(c)
		return
	}

	/** 调用商品服务*/

	tx := dbutil.DB.Begin()
	defer tx.Commit()

	/** 创建不可见订单 */
	orderId, err := SaveNoConfirmOrder(&confirmOrderReq, tx)
	if err != nil {
		logger.Error(err)
		tx.Rollback()
		commonUtils.CreateErrorWithMsg(c, err.Error())
		return
	}

	/** 设置过期时间 */
	err = redis.CacheOrderExpired(orderId)
	if err != nil {
		logger.Error(err)
		tx.Rollback()
		logger.Error("order = %s set cache error ", orderId)
		commonUtils.CreateErrorWithMsg(c, "set cache error")
		return
	}

	commonUtils.CreateSuccess(c, orderId)
}

func (orderService *orderInfoService) FindOrderInfoList(c *gin.Context) {

	var err error

	var orderInfoReq protocol2.OrderInfoReq
	if err = c.ShouldBindJSON(&orderInfoReq); err != nil {
		logger.Error(err.Error())
		commonUtils.CreateError(c)
		return
	}

	logger.Debug(orderInfoReq)

	userId := c.GetString("userId")
	if userId == "" {
		commonUtils.CreateErrorParams(c)
		return
	}

	offset := (orderInfoReq.PageNum - 1) * orderInfoReq.PageSize
	total, err := dao.FindOrderInfoCounts(orderInfoReq.DataType, userId)
	if err != nil {
		logger.Error(err.Error())
		commonUtils.CreateError(c)
		return
	}

	if total < offset {
		commonUtils.CreateNotContent(c)
		return
	}

	orderInfList, err := dao.FindOrderInfoList(offset, orderInfoReq.PageSize, orderInfoReq.DataType, userId)
	if err != nil || len(orderInfList) == 0 {
		commonUtils.CreateNotContent(c)
		return
	}

	orderInfoRespList := make([]protocol2.OrderInfoResp, 0)
	for _, orderInfo := range orderInfList {

		orderInfoResp := protocol2.OrderInfoResp{}
		orderInfoResp.OrderId = orderInfo.OrderId
		orderInfoResp.UserId = orderInfo.UserId
		orderInfoResp.Username = orderInfo.Username
		orderInfoResp.MerchantId = orderInfo.MerchantId
		orderInfoResp.OrderAddress = orderInfo.OrderAddress
		orderInfoResp.IsPostFee = orderInfo.IsPostFee
		orderInfoResp.PostFee = orderInfo.PostFee
		orderInfoResp.GoodsCount = orderInfo.GoodsCount
		orderInfoResp.TotalAmount = orderInfo.TotalAmount
		orderInfoResp.ReallyAmount = orderInfo.ReallyAmount
		orderInfoResp.OrderIdentifier = orderInfo.OrderIdentifier
		orderInfoResp.PayType = orderInfo.PayType
		orderInfoResp.OrderStatus = orderInfo.OrderStatus
		orderInfoResp.ShippingStatus = orderInfo.ShippingStatus
		orderInfoResp.PayStatus = orderInfo.PayStatus
		orderInfoResp.IsPostFee = orderInfo.IsPostFee
		orderInfoResp.Status = ReturnStatus(orderInfo)

		orderInfoResp.OrderGoodsRespList, err = dao.FindOrderGoodsByOrderId(orderInfo.OrderId)
		if err != nil {
			logger.Error(err)
			return
		}

		orderInfoRespList = append(orderInfoRespList, orderInfoResp)
	}

	commonUtils.CreateSuccessByList(c, total, orderInfoRespList)
}

func (orderInfoService *orderInfoService) CancelOrder(c *gin.Context) {

	orderInfoService.Mutex.Lock()
	defer orderInfoService.Mutex.Unlock()

	userId := c.GetString("userId")
	orderId := c.Param("orderId")

	if userId == "" || orderId == "" {
		commonUtils.CreateErrorParams(c)
		return
	}

	err := dao.CancelOrder(orderId, userId)
	if err != nil {
		commonUtils.CreateError(c)
		return
	}

	commonUtils.CreateSuccess(c, nil)
}

func (orderInfoService *orderInfoService) ReceiveIPN(c *gin.Context) {
	orderInfoService.Mutex.Lock()
	defer orderInfoService.Mutex.Unlock()

	logger.Info(c.GetPostForm("verify_sign"))

	requests := c.Request.PostForm
	if requests == nil {
		logger.Error(requests)
		return
	}

	var keys []string
	for k, _ := range requests {
		if k != "verify_sign" {
			keys = append(keys, k)
		}
	}

	if len(keys) == 0 {
		logger.Info("params of null")
		return
	}

	sort.Strings(keys)
	logger.Info("\n request =======> ", requests)

	var str = ""
	params := make(map[string]string)
	for _, v := range keys {
		val := requests[v][0]
		if val != "null" {
			fmt.Printf("%s=%s;/n", v, string(val))
			str = str + string(v) + "=" + string(val) + "&"
			params[string(v)] = string(val)
		}
	}

	mdStr := strings.ToLower(MD5(token))
	mdStr = str + mdStr
	logger.Info("\nprams and token sign string:" + mdStr)

	sign := strings.ToLower(MD5(mdStr))
	logger.Info("\nsign string:" + sign)

	verify_sign := requests["verify_sign"][0]
	params["verify_sign"] = string(verify_sign)

	logger.Infof("params: %s", params)

	if sign == verify_sign {
		logger.Info(params["status"])
		if params["status"] == "success" {
			orderId := params["reference"]
			noteJson := params["note"]

			var note note
			if err := json.Unmarshal([]byte(noteJson), &note); err != nil {
				logger.Error(err)
				return
			}

			message := &message{OrderId: orderId, Note: note}
			logger.Info(message)

			if orderId != "" {
				/** 更改订单状态  */
				if message.Note.Type == "0" {
					//商品
					go func() {

						err := redis.DelOrderExpired(orderId)
						if err != nil {
							logger.Error(err)
						}

						err = dbutil.DB.Exec("update order_info set payment_time = now(), pay_status = ?, pay_type = ? where order_id = ?",
							commonUtils.PAYSTATUS_PAIED, message.Note.PayType, orderId).Error
						if err != nil {
							logger.Error(err)
							return
						}

						var orderInfo model2.OrderInfo
						err = dbutil.DB.Exec("select * from order_info where order_id = ? ", orderId).Find(&orderInfo).Error
						if err != nil {
							logger.Error(err)
							return
						}

						/** 加本人和推荐人积分和佣金 */
						orderInfoService.addIntegral(&orderInfo)

					}()

				}else if message.Note.Type == "1" {
					go func() {
						//服务支付
						var payment model.ServicePaymentSetting
						id := message.Note.SelectId
						if err := dbutil.DB.Exec("select * from service_payment_setting where id = ?", id).
							Find(&payment).Error; err != nil{
							logger.Error(err)
						}else {
							servicePay := model.ServicePayFlowProtocol{
								ServiceId: orderId,
								ExTime: payment.Time,
								Price: payment.Price,
							}
							logic.ServicePayFlowService.UpdateStatus(&servicePay)
						}
					}()
				}
			}
		}

		return
	}
	logger.Info("false")
}

func (orderInfoService *orderInfoService)addIntegral(orderInfo *model2.OrderInfo){

	var err error

	userId := orderInfo.UserId
	totalIntegral := orderInfo.TotalIntegral
	orderId := orderInfo.OrderId

	var userflow = model3.UserIntegralFlow{}
	userflow.UserId = userId
	userflow.Integral = totalIntegral
	userflow.CreateAt = time.Now()
	userflow.UpdateAt = time.Now()
	userflow.OrderId = orderInfo.OrderId

	err = UserDao.UpdateIntegralByUserId(userId, totalIntegral)
	if err != nil {
		logger.Error(err)
		return
	}else {
		err = UserDao.UpdateIntegralByUserId(userId, totalIntegral)
		if err != nil {
			logger.Error(err)
			return
		}

		err = UserDao.SaveUserIntegralFlow(userflow)
		if err != nil {
			logger.Error(err)
		}
	}

	addRecommendIntegral(orderId, userId, totalIntegral / 2)
}

/**
	recommendUserId 被推荐人
 */
func addRecommendIntegral(orderId string, recommendUserId string, totalIntegral int)  {
	recommendUser, err := UserDao.GetRecommendUserIdByUserId(recommendUserId)
	logger.Info(recommendUser)
	if err != nil || recommendUser.RecommendUserId == ""{
		logger.Error(err)
	}else {
		err = UserDao.UpdateIntegralByUserId(recommendUser.RecommendUserId, totalIntegral)
		if err != nil {
			logger.Error(err)
		}

		var userflow = model3.UserIntegralFlow{}
		userflow.CreateAt = time.Now()
		userflow.UpdateAt = time.Now()
		userflow.OrderId = orderId
		userflow.RUserId = recommendUserId
		userflow.UserId = recommendUser.UserId
		userflow.Integral = totalIntegral

		err = UserDao.SaveUserIntegralFlow(userflow)
		if err != nil {
			logger.Error(err)
		}
	}

}


func (orderInfoService *orderInfoService) startUpdateOrderStatus() {
	for {
		select {
		case message := <-worker:
			logger.Info(message)
			orderId := message.OrderId

			if orderId != "" {
				/** 更改订单状态  加本人和推荐人积分和佣金*/
				if message.Note.Type == "0" {
					//商品
					err := dbutil.DB.Exec("update order_info set payment_time = now(), pay_status = ?, pay_type = ? where order_id = ?",
						commonUtils.PAYSTATUS_PAIED, message.Note.PayType, orderId).Error
					if err != nil {
						logger.Error(err)
					}

				}else if message.Note.Type == "1" {
					//服务支付
					var payment model.ServicePaymentSetting
					id := message.Note.SelectId
					if err := dbutil.DB.Exec("select * from service_pay_setting where id = ?", id).Find(&payment).Error; err != nil{
						logger.Error(err)
					}else {
						servicePay := model.ServicePayFlowProtocol{
							ServiceId: orderId,
							ExTime: payment.Time,
							Price: payment.Price,
						}
						logic.ServicePayFlowService.UpdateStatus(&servicePay)
					}
				}
			}
		}
	}
}

func (orderInfoService *orderInfoService) startSaveOrderflow() {
	for {
		select {
	case orderFlow := <-saveWorker:
		logger.Info(orderFlow)
		if orderFlow != nil {
			orderInfoService.SetContext()
			err := orderInfoService.c.Insert(&orderFlow)
			if err != nil {
				logger.Error(err)
			}
			orderInfoService.context.Close()
		}
	}
}
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
