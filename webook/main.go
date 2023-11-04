package main

import (
	"net/http"
	"strings"
	"time"

	"go_learning/config"
	"go_learning/internal/repository"
	"go_learning/internal/repository/dao"
	"go_learning/internal/service"
	"go_learning/internal/web"
	"go_learning/internal/web/middleware"
	ratelimit "go_learning/pkg/ginx/middleware/retelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()

	server := initWebServer()
	initUserHdl(db, server)
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	server.Run(":8081")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	// "root:root@tcp(localhost:13316)/webook"
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,

		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
		//AllowMethods: []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的 Middleware")
	})

	// "localhost:6379"
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 20).Build())

	// useSession(server)
	useJWT(server)

	return server
}

func useJWT(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
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
