package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserLogin struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type UserWithRowNumber struct {
	RowNumber   int
	ID          int
	Name        string
	Username    string
	Email       string
	PhoneNumber string
	RoleName    string
}

type Account struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
	Email       string `form:"email"`
	Name        string `form:"name"`
	PhoneNumber string `form:"phonenumber"`
	RoleID      int    `form:"role_id"`
}

func HashPassword(user *User, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func CheckPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// BeforeCreate là hook GORM, chạy trước khi tạo một bản ghi mới.
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err := HashPassword(user, user.Password); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate là hook GORM, chạy trước khi cập nhật một bản ghi.
func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if err := HashPassword(user, user.Password); err != nil {
		return err
	}
	return nil
}
