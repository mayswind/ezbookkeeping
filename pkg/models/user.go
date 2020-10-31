package models

type UserType byte

const (
	USER_TYPE_NORMAL      UserType = 0
	USER_TYPE_ADMIN       UserType = 63
	USER_TYPE_SUPER_ADMIN UserType = 127
)

type User struct {
	Uid               int64    `xorm:"PK"`
	Username          string   `xorm:"VARCHAR(32) UNIQUE NOT NULL"`
	Email             string   `xorm:"VARCHAR(100) UNIQUE NOT NULL"`
	Nickname          string   `xorm:"VARCHAR(64) NOT NULL"`
	Password          string   `xorm:"VARCHAR(64) NOT NULL"`
	Salt              string   `xorm:"VARCHAR(10) NOT NULL"`
	Rands             string   `xorm:"VARCHAR(10) NOT NULL"`
	Type              UserType `xorm:"TINYINT NOT NULL"`
	IsAdmin           bool     `xorm:"NOT NULL"`
	Deleted           bool     `xorm:"NOT NULL"`
	EmailVerified     bool     `xorm:"NOT NULL"`
	CreatedUnixTime   int64
	UpdatedUnixTime   int64
	DeletedUnixTime   int64
	LastLoginUnixTime int64
}

type UserLoginRequest struct {
	LoginName string `json:"loginName" binding:"required,notBlank,max=100,validUsername|validEmail"`
	Password  string `json:"password" binding:"required,min=6,max=128"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,notBlank,max=32,validUsername"`
	Email    string `json:"email" binding:"required,notBlank,max=100,validEmail"`
	Nickname string `json:"nickname" binding:"required,notBlank,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type UserProfileUpdateRequest struct {
	Email       string `json:"email" binding:"omitempty,notBlank,max=100,validEmail"`
	Nickname    string `json:"nickname" binding:"omitempty,notBlank,max=64"`
	Password    string `json:"password" binding:"omitempty,min=6,max=128"`
	OldPassword string `json:"oldPassword" binding:"omitempty,min=6,max=128"`
}

type UserProfileResponse struct {
	Uid         string   `json:"uid"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Nickname    string   `json:"nickname"`
	Type        UserType `json:"type"`
	CreatedAt   int64    `json:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt"`
	LastLoginAt int64    `json:"lastLoginAt"`
}
