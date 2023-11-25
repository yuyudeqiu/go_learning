package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go_learning/internal/domain"
	"go_learning/internal/repository/cache"
	cachemocks "go_learning/internal/repository/cache/mocks"
	"go_learning/internal/repository/dao"
	daomocks "go_learning/internal/repository/dao/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO)

		ctx context.Context
		uid int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功,缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExists)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id:          uid,
						Email:       sql.NullString{String: "123@qq.com", Valid: true},
						Password:    "123456",
						Birthday:    100,
						Description: "about me",
						Phone:       sql.NullString{String: "114514", Valid: true},
						Ctime:       nowMs,
						Utime:       102,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:          123,
					Email:       "123@qq.com",
					Password:    "123456",
					Birthday:    time.UnixMilli(100),
					Description: "about me",
					Phone:       "114514",
					Ctime:       now,
				}).Return(nil)
				return c, d
			},

			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:          123,
				Email:       "123@qq.com",
				Password:    "123456",
				Birthday:    time.UnixMilli(100),
				Description: "about me",
				Phone:       "114514",
				Ctime:       now,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, ud := tc.mock(ctrl)
			svc := NewCacheUserRepository(ud, uc)
			user, err := svc.FindById(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
