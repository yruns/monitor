package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Avatar         string
	Money          string
}

const (
	PasswordCost = 12 // 密码加密难度
	Active       = "1"
)

func (u *User) SetPassword(pwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), PasswordCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

func (u *User) CheckPassword(password, hashedPasswordFromDB string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
