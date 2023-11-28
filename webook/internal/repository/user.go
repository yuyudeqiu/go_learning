package repository

import (
	"context"
	"database/sql"
	"time"

	"go_learning/internal/domain"
	"go_learning/internal/repository/cache"
	"go_learning/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
	FindByWechat(ctx context.Context, openId string) (user domain.User, err error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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

func (repo *CacheUserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *CacheUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		NickName:    u.NickName,
		Birthday:    time.UnixMilli(u.Birthday),
		Email:       u.Email.String,
		Phone:       u.Phone.String,
		Password:    u.Password,
		Description: u.Description,
		Ctime:       time.UnixMilli(u.Ctime),
		WechatInfo: domain.WechatInfo{
			UnionId: u.WechatUnionId.String,
			OpenId:  u.WechatOpenId.String,
		},
	}
}

func (repo *CacheUserRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:            user.Id,
		NickName:      user.NickName,
		Birthday:      user.Birthday.UnixMilli(),
		Email:         sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:         sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Description:   user.Description,
		Ctime:         user.Ctime.UnixMilli(),
		WechatOpenId:  sql.NullString{String: user.WechatInfo.OpenId, Valid: user.WechatInfo.OpenId != ""},
		WechatUnionId: sql.NullString{String: user.WechatInfo.UnionId, Valid: user.WechatInfo.UnionId != ""},
	}
}

func (repo *CacheUserRepository) FindByWechat(ctx context.Context, openId string) (user domain.User, err error) {
	ue, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(ue), nil
}
