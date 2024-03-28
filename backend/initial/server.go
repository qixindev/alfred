package initial

import (
	"alfred/backend"
	"alfred/backend/pkg/cache"
	"alfred/backend/pkg/config/env"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"time"
)

func initSystem() error {
	var err error
	if err = InitConfig(); err != nil { // 初始化配置
		return err
	}

	// 初始化日志
	global.LOG = Zap()
	zap.ReplaceGlobals(global.LOG)
	if global.LOG == nil {
		return errors.New("init zap log err")
	}

	if err = InitDB(); err != nil {
		return errors.WithMessage(err, "InitDB err")
	}

	global.CodeCache, err = cache.NewBigCache(5 * 60 * time.Second)
	if err != nil {
		return errors.WithMessage(err, "init code big cache err")
	}

	global.StateCache, err = cache.NewBigCache(5 * 60 * time.Second)
	if err != nil {
		return errors.WithMessage(err, "init state big cache err")
	}

	if env.GetReleaseType() == "first" {
		migrateList := getMigrateModel()
		if err = global.DB.AutoMigrate(migrateList...); err != nil {
			return errors.WithMessage(err, "migrate db err")
		}
		if err = InitDefaultTenant(); err != nil {
			return errors.WithMessage(err, "InitDefaultTenant err")
		}
	}

	return nil
}

func StartServer(_ *cobra.Command, _ []string) {
	var err error
	for i := 0; i < 10; i++ {
		if err = initSystem(); err == nil {
			break
		}
		fmt.Println("init system err:", err)
		time.Sleep(time.Second * 30)
	}

	if err != nil {
		os.Exit(1)
		return
	}

	fmt.Println(fmt.Sprintf("### %s: %v\n", env.GetDeployType(), utils.StructToString(global.CONFIG)))

	r := gin.New()
	cookieSecret := GetSessionSecret()
	store := cookie.NewStore(cookieSecret)
	store.Options(sessions.Options{
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	})
	r.Use(sessions.Sessions("QixinAuth", store), gin.Logger())
	backend.AddRoutes(r)

	if err = r.Run(":80"); err != nil {
		fmt.Println("server run err: ", err)
		return
	}
}
