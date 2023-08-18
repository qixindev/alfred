package cmd

import (
	"accounts/initial"
	"accounts/internal"
	"accounts/pkg/config/env"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func initSystem() error {
	var err error
	if err = initial.InitConfig(); err != nil { // 初始化配置
		return err
	}

	// 初始化日志
	global.LOG = initial.Zap()
	zap.ReplaceGlobals(global.LOG)
	if global.LOG == nil {
		return errors.New("init zap log err")
	}

	if err = initial.InitDB(); err != nil {
		return err
	}

	return nil
}

func startServer() {
	var err error
	for i := 0; i < 10; i++ {
		if err = initSystem(); err == nil {
			break
		}
		time.Sleep(time.Second * 30)
	}

	if err != nil {
		fmt.Println("init system err:", err)
		return
	}

	fmt.Println(fmt.Sprintf("### %s: %v\n", env.GetDeployType(), utils.StructToString(global.CONFIG)))

	r := gin.Default()
	r.Use(cors.Default())
	cookieSecret := initial.GetSessionSecret()
	store := cookie.NewStore(cookieSecret)
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	})
	r.Use(sessions.Sessions("QixinAuth", store))
	internal.AddRoutes(r)

	if err = r.Run(":80"); err != nil {
		fmt.Println("server run err: ", err)
		return
	}
}
