package auth

import (
	"github.com/gin-gonic/gin"
)

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

type AuthProvider interface {
	// Auth Get to external auth. Return redirect location.
	Auth(string) string

	// Login Callback when auth completed.
	Login(*gin.Context) (*UserInfo, error)
}
