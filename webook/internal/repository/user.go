package repository

import (
	"context"
	"database/sql"
	"time"

	"go_learning/internal/domain"
	"go_learning/internal/repository/cache"
	"go_learning/internal/repository/dao"

	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, id)
	if err == nil {
		return du, nil
	}
	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	repo.cache.Set(ctx, du)
	return du, nil
}

func (repo *UserRepository) UpdateNonZeroFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		NickName:    u.NickName,
		Birthday:    time.UnixMilli(u.Birthday),
		Email:       u.Email.String,
		Phone:       u.Phone.String,
		Password:    u.Password,
		Description: u.Description,
	}
}

func (repo *UserRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:          user.Id,
		NickName:    user.NickName,
		Birthday:    user.Birthday.UnixMilli(),
		Email:       sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:       sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Description: user.Description,
	}
}
