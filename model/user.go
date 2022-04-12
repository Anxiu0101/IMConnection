package model

import (
	"MedicalCare/conf"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	// Account Info
	UserName string `json:"username" gorm:"column:username;not null;uniqueIndex"`
	Password string `json:"password" gorm:"column:password"`
	Avatars  string `json:"avatars" gorm:"column:avatars"`

	// User Info
	Email  string `json:"email" gorm:"type:varchar(100);unique"`
	Gender int    `json:"gender" gorm:"size:3"`
	Age    int    `json:"age" gorm:"size:8"`
	Tel    int    `json:"tel" gorm:"size:13"`

	State bool `json:"state" gorm:"column:state"`
}

//SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), conf.AppSetting.PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

//CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

/* serialization */

// UserInfo 用户资料结构体
type UserInfo struct {
	Email  string `json:"email" form:"email"`
	Gender int    `json:"gender" form:"gender"`
	Age    int    `json:"age" form:"age"`
	Tel    int    `json:"tel" form:"tel"`
}

type AccountInfo struct {
	ID       uint   `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	Avatars  string `json:"avatars" form:"avatars"`
}

func BuildUserInfo(user User) UserInfo {
	return UserInfo{
		Email:  user.Email,
		Gender: user.Gender,
		Age:    user.Age,
		Tel:    user.Tel,
	}
}

func BuildAccountInfo(user User) AccountInfo {
	return AccountInfo{
		ID:       user.ID,
		UserName: user.UserName,
		Avatars:  user.Avatars,
	}
}

func BuildUser(user User) (AccountInfo, UserInfo) {
	return AccountInfo{
			ID:       user.ID,
			UserName: user.UserName,
			Avatars:  user.Avatars,
		}, UserInfo{
			Email:  user.Email,
			Gender: user.Gender,
			Age:    user.Age,
			Tel:    user.Tel,
		}
}
