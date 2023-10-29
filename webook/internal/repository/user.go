package repository

import (
	"context"
	"time"

	"go_learning/internal/domain"
	"go_learning/internal/repository/dao"

	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		NickName:    u.NickName,
		Birthday:    time.UnixMilli(u.Birthday),
		Email:       u.Email,
		Password:    u.Password,
		Description: u.Description,
	}
}

func (repo *UserRepository) UpdateNonZeroFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *UserRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:          user.Id,
		NickName:    user.NickName,
		Birthday:    user.Birthday.UnixMilli(),
		Email:       user.Email,
		Description: user.Description,
	}
}
