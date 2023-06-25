package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Filter[S any, T any](s []S, f func(S) T) []T {
	var l []T
	for _, i := range s {
		l = append(l, f(i))
	}
	return l
}

func GetHostWithScheme(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if s := c.Request.Header.Get("X-Forwarded-Proto"); s != "" {
		scheme = s
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

func GetString(v interface{}) string {
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// DeferErr 处理defer返回的错误警告
func DeferErr(errFunc func() error) {
	if err := errFunc(); err != nil {
		fmt.Println("### Defer err: ", err)
	}
}
