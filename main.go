package main

import (
	"accounts/data"
	"accounts/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
)

var secret = []byte("secret")

func main() {
	r := gin.Default()
	r.Use(sessions.Sessions("QixinAuth", cookie.NewStore(secret)))
	routes.AddRoutes(&r.RouterGroup)
	if err := data.InitDB(); err != nil {
		log.Fatal(err)
		return
	}
	if err := r.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
