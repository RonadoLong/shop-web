package main

import (
	"flag"
	"io"
	"math/rand"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"shop-web/common/cron"
	"shop-web/common/dbutil"
	"shop-web/common/log"
	"shop-web/common/redis"
	"shop-web/conf"
	"shop-web/module/engine"
	"syscall"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/HaroldHoo/id_generator"
	"github.com/gin-gonic/gin"
	"shop-web/module/order/logic"
)

var Logger = log.NewLogger(os.Stdout)

func init() {

	enviroment := flag.String("e", "pro", "")
	Logger.Info(enviroment)

	flag.Parse()
	conf.Init(*enviroment)
	c := conf.GetConfig()
	Logger.Info(c.GetString("mysql.host"))

	if c.GetString("config") == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//cron
	cron.InitCron(c.GetString("stockCron"))

	//mysql
	dbutil.ConnectDB(c.GetString("mysql.host"), c.GetString("mysql.logLevel"))

	//redis
	redis.InitRedis(c.GetString("redis.host"), c.GetString("redis.password"))

	id_generator.DefaultInstanceId = 154

	rand.Seed(time.Now().Unix())

	gin.DefaultWriter = io.MultiWriter(os.Stdout)

	//logic.OrderInfoService.InitDetail()
}

func main() {
	c := conf.GetConfig()
	router := engine.ClientEngine()
	router.Use(gin.Logger())
	server := &http.Server{
		Addr:           c.GetString("server.host"),
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ginpprof.Wrap(router)
	handleSignal(server)

	Logger.Infof("shop (v%s) is running [%s]", c.GetString("server.version"), c.GetString("server.host"))
	server.ListenAndServe()

	logic.OrderInfoService.InitDetail()
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		s := <-c
		Logger.Infof("got signal [%s], exiting shop now", s)
		if err := server.Close(); nil != err {
			Logger.Errorf("server close failed: ", err)
		}
		dbutil.DisconnectDB()
		Logger.Infof("shop exited")
		os.Exit(0)
	}()
}
