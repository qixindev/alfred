package main

import (
	"accounts/data"
	"accounts/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
)

var secret = []byte("secret")

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(sessions.Sessions("QixinAuth", cookie.NewStore(secret)))
	routes.AddRoutes(r)
	if err := data.InitDB(); err != nil {
		log.Fatal(err)
		return
	}
	if err := r.Run(":8086"); err != nil {
		log.Fatal(err)
		return
	}
}
