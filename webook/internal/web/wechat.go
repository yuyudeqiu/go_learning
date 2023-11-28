package web

import (
	"go_learning/internal/service"
	"go_learning/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService

	key             []byte
	stateCookieName string
}
