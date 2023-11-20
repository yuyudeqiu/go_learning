package ioc

import (
	"go_learning/internal/service/localsms"
	"go_learning/internal/service/sms"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
}
