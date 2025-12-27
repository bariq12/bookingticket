package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"github.com/bariq12/bookingticket/models"
	"github.com/bariq12/bookingticket/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type AuthService struct{
	repository models.AuthRepository
}

func (s *AuthService) Login(ctx context.Context, loginData *models.AuthCrendetial) (string, *models.User, error){
	user, err := s.repository.GetUser(ctx, "email = ?", loginData.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return "", nil, fmt.Errorf("error credentials")
		}
		return "", nil, err
	}
	if models.MatchesHash(loginData.Password, user.Password){
		return "", nil, fmt.Errorf("error credentials")
	}
	claims := jwt.MapClaims{
		"id" : user.ID,
		"role" : user.Role,
		"exp" : time.Now().Add(time.Hour * 168).Unix(),
	
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_Secret"))
	if err != nil {
		return "", nil, fmt.Errorf("error generating token: %v", err)
	}

	return token, user, nil
}

func (s *AuthService) Register(ctx context.Context, registerData *models.AuthCrendetial) (string, *models.User, error){
	if !models.Isvaliduseremail(registerData.Email) {
		return "", nil, fmt.Errorf("invalid email format")
	}
	if _, err := s.repository.GetUser(ctx, "email = ?", registerData.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, fmt.Errorf("email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("error hashing password: %v", err)
	}
	registerData.Password = string(hashedPassword)
	user, err := s.repository.RegisterUser(ctx, registerData)
	if err != nil {
		return "", nil, fmt.Errorf("error registering user: %v", err)
	}
	claims := jwt.MapClaims{
		"id" : user.ID,
		"role" : user.Role,
		"exp" : time.Now().Add(time.Hour * 168).Unix(),
	
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_Secret"))
	if err != nil {
		return "", nil, fmt.Errorf("error generating token: %v", err)
	}

	return token, user, nil
}
func NewAuthService(repository models.AuthRepository) models.AuthService {
	return &AuthService{
		repository : repository,
	}
}