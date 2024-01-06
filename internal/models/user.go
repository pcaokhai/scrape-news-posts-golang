package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID			string		`json:"userId" db:"user_id" redis:"user_id" validate:"omitempty"`
	Username		string		`json:"username" db:"username" redis:"username" validate:"required"`
	Password 		string		`json:"password,omitempty" db:"password" redis:"password" validate:"omitempty,required"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}


func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// return string(hashedPassword) == u.Password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}
