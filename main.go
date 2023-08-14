package main

import (
	"accounts/internal/controller"
	"accounts/internal/global"
	initial2 "accounts/internal/initial"
	"accounts/pkg/config/env"
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
	if err = initial2.InitConfig(); err != nil { // 初始化配置
		fmt.Println("Init Config error: " + err.Error())
		return err
	}

	// 初始化日志
	global.LOG = initial2.Zap()
	zap.ReplaceGlobals(global.LOG)
	if global.LOG == nil {
		fmt.Println("init zap log err: zap log is nil")
		return errors.New("init zap log err")
	}

	if err = initial2.InitDB(); err != nil {
		fmt.Println("Init DB error: ", err)
		return err
	}

	if err = initial2.CheckFirstRun(); err != nil {
		fmt.Println("first run err: ", err)
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
	cookieSecret := initial2.GetSessionSecret()
	store := cookie.NewStore(cookieSecret)
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	})
	r.Use(sessions.Sessions("QixinAuth", store))
	controller.AddRoutes(r)

	if err = r.Run(":80"); err != nil {
		fmt.Println("server run err: ", err)
		return
	}
}
