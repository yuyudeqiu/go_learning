package ratelimit

import (
	"context"
	"errors"
	"testing"

	"go_learning/internal/service/sms"
	smsmocks "go_learning/internal/service/sms/mocks"
	"go_learning/pkg/limiter"
	limitermocks "go_learning/pkg/limiter/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		wantErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return svc, l
			},
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(true, nil)
				return svc, l
			},
			wantErr: errLimited,
		},
		{
			name: "限流错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(true, errors.New("限流错误"))
				return svc, l
			},
			wantErr: errors.New("限流错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)
			svc := NewRateLimitSMSService(smsSvc, l)
			err := svc.Send(context.Background(), "abc", []string{"123"}, "123456")
			assert.Equal(t, tc.wantErr, err)
		})
	}

}
