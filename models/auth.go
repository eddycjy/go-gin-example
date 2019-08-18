package models

import (
	"github.com/jinzhu/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CheckAuth checks if authentication information exists
func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}

//注册用户
func AddAuth(username, password string) (bool, error) {
	err := db.Create(&Auth{Username: username, Password: password}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

//修改密码
func ResetPassword(username, password, newPassword string) (bool, error) {
	err := db.Model(Auth{}).Where(Auth{Username: username, Password: password}).Update("password", newPassword).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

func Detail(username string) (Auth, error) {
	var auth Auth
	err := db.First(&auth).Where("username", username).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	return auth, nil
}

func All() ([]Auth, error) {
	var auth []Auth
	err := db.Find(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	return auth, nil
}
