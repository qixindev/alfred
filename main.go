package main

import (
	"accounts/config/env"
	"accounts/global"
	"accounts/initial"
	"accounts/server"
	"accounts/utils"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func InitSystem() error {
	var err error
	if err = initial.InitConfig(); err != nil { // 初始化配置
		fmt.Println("Init Config error: " + err.Error())
		return err
	}

	// 初始化日志
	global.LOG = initial.Zap()
	zap.ReplaceGlobals(global.LOG)
	if global.LOG == nil {
		fmt.Println("init zap log err: zap log is nil")
		return errors.New("init zap log err")
	}

	if err = initial.InitDB(); err != nil {
		fmt.Println("Init DB error: ", err)
		return err
	}

	return nil
}

func main() {
	var err error
	for i := 0; i < 10; i++ {
		if err = InitSystem(); err == nil {
			break
		}
		time.Sleep(time.Second * 30)
	}

	if err != nil {
		return
	}

	global.LOG.Info(fmt.Sprintf("### %s: %v\n", env.GetDeployType(), utils.StructToString(global.CONFIG)))

	r := gin.Default()
	r.Use(cors.Default())
	cookieSecret := initial.GetSessionSecret()
	r.Use(sessions.Sessions("QixinAuth", cookie.NewStore(cookieSecret)))
	server.AddRoutes(r)

	if err = r.Run(":8086"); err != nil {
		fmt.Println("server run err: ", err)
		return
	}
}
