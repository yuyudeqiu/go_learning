//go:build wireinject

package main

import (
	"go_learning/internal/repository"
	"go_learning/internal/repository/cache"
	"go_learning/internal/repository/dao"
	"go_learning/internal/service"
	"go_learning/internal/web"
	ijwt "go_learning/internal/web/jwt"
	"go_learning/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLocalCache,

		dao.NewUserDAO,

		//cache.NewCodeCache,
		cache.NewLocalCodeCache,
		cache.NewUserCache,

		repository.NewCacheUserRepository,
		repository.NewCodeRepository,

		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,

		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebService,
	)
	return gin.Default()
}
