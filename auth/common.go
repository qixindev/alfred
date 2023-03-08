package auth

import "github.com/gin-gonic/gin"

type UserInfo struct {
	Sub         string
	Name        string
	FirstName   string
	LastName    string
	DisplayName string
	Email       string
	Phone       string
	Picture     string
}

func GetString(v interface{}) string {
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type AuthProvider interface {
	// Auth Get to external auth. Return redirect location.
	Auth(string) string

	// Login Callback when auth completed.
	Login(*gin.Context) (*UserInfo, error)
}
