package models

import (
	"user-service/src/api/generated"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// The reason for this model is that in a perfect world I would have
	// implemented UUID and timestamp types here.
	ID        string `json:"id"` // Saving internally as string because that's the type we're using for our gRPC message.
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Use bcrypt to salt and hash the password.
	return string(bytes)
}

func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// I could do a normal go type cast for this
// but in a perfect world where I properly handle
// UUIDs and timestamps this would hold real logic
// so I'm keeping it as a point to make later.
func (u *User) ToRPC() *generated.User {
	if u == nil {
		return nil
	}

	return &generated.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Password:  u.Password,
		Email:     u.Email,
		Country:   u.Country,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromRPC(u *generated.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:        u.Id,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Nickname:  u.Nickname,
		Password:  u.Password,
		Email:     u.Email,
		Country:   u.Country,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
