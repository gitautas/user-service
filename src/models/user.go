package models

import (
	"time"
	"user-service/src/api/generated"

	"golang.org/x/crypto/bcrypt"
)

type User generated.User

func (u *User) HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Use bcrypt to salt and hash the password.
	return string(bytes)
}

func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) GenerateTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
