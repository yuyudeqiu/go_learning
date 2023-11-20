package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	server.Run(":8081")
}

//func useSession(server *gin.Engine) {
//	login := middleware.LoginMiddlewareBuilder{}
//	store, err := redis.NewStore(16, "tcp",
//		"localhost:6379", "",
//		[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
//		[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
//	if err != nil {
//		panic(err)
//	}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
