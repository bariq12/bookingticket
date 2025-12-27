package models

import (
	"context"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type AuthCrendetial struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, registerData *AuthCrendetial) (*User, error)
	GetUser(ctx context.Context, query interface{}, args ...interface{}) (*User, error)
}

type AuthService interface{
	Login(ctx context.Context, loginData *AuthCrendetial) (string, *User, error)
	Register(ctx context.Context, registerData *AuthCrendetial) (string, *User, error)
}

func MatchesHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return  err == nil
}

func Isvaliduseremail(email string) bool {
	_, err := mail.ParseAddress(email)
	return  err == nil
} 