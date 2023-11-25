package service

import (
	"errors"
	"testing"

	"go_learning/internal/domain"
	"go_learning/internal/repository"
	repomocks "go_learning/internal/repository/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("hello@ccC56")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("hello@ccC56"))
	assert.NoError(t, err)
}

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 预期输入
		ctx      *gin.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$TLj9IGIGnPrFjpVnsEhWhuXysM8.9DTbheJ9EJY5H4loo2ozqWTNi",
						Phone:    "123456",
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello@ccC56",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$TLj9IGIGnPrFjpVnsEhWhuXysM8.9DTbheJ9EJY5H4loo2ozqWTNi",
				Phone:    "123456",
			},
			wantErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "$2a$10$TLj9IGIGnPrFjpVnsEhWhuXysM8.9DTbheJ9EJY5H4loo2ozqWTNi",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$TLj9IGIGnPrFjpVnsEhWhuXysM8.9DTbheJ9EJY5H4loo2ozqWTNi",
						Phone:    "123456",
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello@ccC55",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},

		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "hello@ccC56",

			wantUser: domain.User{},
			wantErr:  errors.New("db错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
