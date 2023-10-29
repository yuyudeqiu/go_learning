package domain

import "time"

type User struct {
	Id          int64
	NickName    string
	Birthday    time.Time
	Email       string
	Password    string
	Description string

	// UTC 0 的时区
	Ctime time.Time
}
