package models

import (
	"time"

	"github.com/mayswind/lab/pkg/utils"
)

// UserType represents user type
type UserType byte

// User types
const (
	USER_TYPE_NORMAL      UserType = 0
	USER_TYPE_ADMIN       UserType = 63
	USER_TYPE_SUPER_ADMIN UserType = 127
)

// WeekDay represents week day
type WeekDay byte

// Week days
const (
	WEEKDAY_SUNDAY    WeekDay = 0
	WEEKDAY_MONDAY    WeekDay = 1
	WEEKDAY_TUESDAY   WeekDay = 2
	WEEKDAY_WEDNESDAY WeekDay = 3
	WEEKDAY_THURSDAY  WeekDay = 4
	WEEKDAY_FRIDAY    WeekDay = 5
	WEEKDAY_SATURDAY  WeekDay = 6
	WEEKDAY_INVALID   WeekDay = 255
)

// TransactionEditScope represents the scope which transaction can be edited
type TransactionEditScope byte

// Historical Transaction Edit Scopes
const (
	TRANSACTION_EDIT_SCOPE_NONE                TransactionEditScope = 0
	TRANSACTION_EDIT_SCOPE_ALL                 TransactionEditScope = 1
	TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER      TransactionEditScope = 2
	TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER   TransactionEditScope = 3
	TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER  TransactionEditScope = 4
	TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER TransactionEditScope = 5
	TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER  TransactionEditScope = 6
	TRANSACTION_EDIT_SCOPE_INVALID             TransactionEditScope = 255
)

// User represents user data stored in database
type User struct {
	Uid                  int64                `xorm:"PK"`
	Username             string               `xorm:"VARCHAR(32) UNIQUE NOT NULL"`
	Email                string               `xorm:"VARCHAR(100) UNIQUE NOT NULL"`
	Nickname             string               `xorm:"VARCHAR(64) NOT NULL"`
	Password             string               `xorm:"VARCHAR(64) NOT NULL"`
	Salt                 string               `xorm:"VARCHAR(10) NOT NULL"`
	Rands                string               `xorm:"VARCHAR(10) NOT NULL"`
	Type                 UserType             `xorm:"TINYINT NOT NULL"`
	DefaultCurrency      string               `xorm:"VARCHAR(3) NOT NULL"`
	FirstDayOfWeek       WeekDay              `xorm:"TINYINT NOT NULL"`
	TransactionEditScope TransactionEditScope `xorm:"TINYINT NOT NULL"`
	IsAdmin              bool                 `xorm:"NOT NULL"`
	Deleted              bool                 `xorm:"NOT NULL"`
	EmailVerified        bool                 `xorm:"NOT NULL"`
	CreatedUnixTime      int64
	UpdatedUnixTime      int64
	DeletedUnixTime      int64
	LastLoginUnixTime    int64
}

// UserBasicInfo represents a view-object of user basic info
type UserBasicInfo struct {
	Username             string               `json:"username"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	DefaultCurrency      string               `json:"defaultCurrency"`
	FirstDayOfWeek       WeekDay              `json:"firstDayOfWeek"`
	TransactionEditScope TransactionEditScope `json:"transactionEditScope"`
}

// UserLoginRequest represents all parameters of user login request
type UserLoginRequest struct {
	LoginName string `json:"loginName" binding:"required,notBlank,max=100,validUsername|validEmail"`
	Password  string `json:"password" binding:"required,min=6,max=128"`
}

// UserRegisterRequest represents all parameters of user registering request
type UserRegisterRequest struct {
	Username        string  `json:"username" binding:"required,notBlank,max=32,validUsername"`
	Email           string  `json:"email" binding:"required,notBlank,max=100,validEmail"`
	Nickname        string  `json:"nickname" binding:"required,notBlank,max=64"`
	Password        string  `json:"password" binding:"required,min=6,max=128"`
	DefaultCurrency string  `json:"defaultCurrency" binding:"required,len=3,validCurrency"`
	FirstDayOfWeek  WeekDay `json:"firstDayOfWeek" binding:"min=0,max=6"`
}

// UserProfileUpdateRequest represents all parameters of user updating profile request
type UserProfileUpdateRequest struct {
	Email                string                `json:"email" binding:"omitempty,notBlank,max=100,validEmail"`
	Nickname             string                `json:"nickname" binding:"omitempty,notBlank,max=64"`
	Password             string                `json:"password" binding:"omitempty,min=6,max=128"`
	OldPassword          string                `json:"oldPassword" binding:"omitempty,min=6,max=128"`
	DefaultCurrency      string                `json:"defaultCurrency" binding:"omitempty,len=3,validCurrency"`
	FirstDayOfWeek       *WeekDay              `json:"firstDayOfWeek" binding:"omitempty,min=0,max=6"`
	TransactionEditScope *TransactionEditScope `json:"transactionEditScope" binding:"omitempty,min=0,max=7"`
}

// UserProfileUpdateResponse represents the data returns to frontend after updating profile
type UserProfileUpdateResponse struct {
	User     *UserBasicInfo `json:"user"`
	NewToken string         `json:"newToken,omitempty"`
}

// UserProfileResponse represents a view-object of user profile
type UserProfileResponse struct {
	Username             string               `json:"username"`
	Email                string               `json:"email"`
	Nickname             string               `json:"nickname"`
	Type                 UserType             `json:"type"`
	DefaultCurrency      string               `json:"defaultCurrency"`
	FirstDayOfWeek       WeekDay              `json:"firstDayOfWeek"`
	TransactionEditScope TransactionEditScope `json:"transactionEditScope"`
	LastLoginAt          int64                `json:"lastLoginAt"`
}

// CanEditTransactionByTransactionTime returns whether this user can edit transaction with specified transaction time
func (u *User) CanEditTransactionByTransactionTime(transactionTime int64, utcOffset int) bool {
	if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_NONE {
		return false
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_ALL {
		return true
	}

	now := time.Now()

	transactionUnixTime := utils.GetUnixTimeFromTransactionTime(transactionTime)

	if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER {
		return transactionUnixTime >= now.Unix()-24*60*60
	}

	_, serverUtcOffset := now.Zone()
	serverTodayFirstUnixTime := now.Unix() - int64(now.Hour()*60*60+now.Minute()*60+now.Second())
	clientTodayFirstUnixTime := serverTodayFirstUnixTime + int64(utcOffset*60-serverUtcOffset)

	if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER {
		return transactionUnixTime >= clientTodayFirstUnixTime
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER {
		dayOfWeek := int(now.Weekday()) - int(u.FirstDayOfWeek)

		if dayOfWeek < 0 {
			dayOfWeek += 7
		}

		clientWeekFirstUnixTime := clientTodayFirstUnixTime - int64(dayOfWeek*24*60*60)
		return transactionUnixTime >= clientWeekFirstUnixTime
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER {
		clientMonthFirstUnixTime := clientTodayFirstUnixTime - int64((now.Day()-1)*24*60*60)
		return transactionUnixTime >= clientMonthFirstUnixTime
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER {
		clientYearFirstUnixTime := clientTodayFirstUnixTime - int64((now.YearDay()-1)*24*60*60)
		return transactionUnixTime >= clientYearFirstUnixTime
	}

	return false
}

// ToUserBasicInfo returns a user basic view-object according to database model
func (u *User) ToUserBasicInfo() *UserBasicInfo {
	return &UserBasicInfo{
		Username:             u.Username,
		Email:                u.Email,
		Nickname:             u.Nickname,
		DefaultCurrency:      u.DefaultCurrency,
		FirstDayOfWeek:       u.FirstDayOfWeek,
		TransactionEditScope: u.TransactionEditScope,
	}
}

// ToUserProfileResponse returns a user profile view-object according to database model
func (u *User) ToUserProfileResponse() *UserProfileResponse {
	return &UserProfileResponse{
		Username:             u.Username,
		Email:                u.Email,
		Nickname:             u.Nickname,
		Type:                 u.Type,
		DefaultCurrency:      u.DefaultCurrency,
		FirstDayOfWeek:       u.FirstDayOfWeek,
		TransactionEditScope: u.TransactionEditScope,
		LastLoginAt:          u.LastLoginUnixTime,
	}
}
