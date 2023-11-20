package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱重复")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx *gin.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx *gin.Context, id int64) (User, error)
	UpdateById(ctx *gin.Context, user User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return nil
}

func (dao *GORMUserDAO) FindByEmail(ctx *gin.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("`phone` = ?", phone).First(&res).Error
	return res, err
}

func (dao *GORMUserDAO) FindById(ctx *gin.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) UpdateById(ctx *gin.Context, user User) error {
	return dao.db.WithContext(ctx).Model(&user).Where("`id` = ?", user.Id).Updates(map[string]any{
		"utime":       time.Now().UnixMilli(),
		"nick_name":   user.NickName,
		"birthday":    user.Birthday,
		"description": user.Description,
	}).Error
}

type User struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	NickName    string `gorm:"type:varchar(128)"`
	Birthday    int64
	Email       sql.NullString `gorm:"unique"`
	Password    string         `gorm:"type:varchar(64)"`
	Description string         `gorm:"type:varchar(512)"`
	Phone       sql.NullString `gorm:"unique"`

	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64

	// json 存储
	//Addr string
}
