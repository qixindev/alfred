package middlewares

import (
	"alfred/backend/pkg/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func GetPanicStackInfo(msg string, err any, skip int) string {
	res := fmt.Sprintf("\n[Recovery] panic recovered: %s\n[Error] %v", msg, err)
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if !strings.Contains(file, "github.com") && !strings.Contains(file, "gorm.io/") &&
			!strings.Contains(file, "net/http") && !strings.Contains(file, "runtime/") {
			res += fmt.Sprintf("\n\t%s:%d %s", file, line, runtime.FuncForPC(pc).Name())
		}
	}
	return res + "\n"
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				req := fmt.Sprintf("method:%s path:%s", c.Request.Method, c.Request.URL.Path)
				_ = c.Error(err.(error))
				global.LOG.Error(GetPanicStackInfo(req, err, 3))
				if brokenPipe {
					c.Abort()
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "server panic"})
			}
		}()
		c.Next()
	}
}
