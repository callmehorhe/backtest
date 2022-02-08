package service

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math/rand"
	"time"

	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	passwordLetters = "0123456789abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt       = "opie435qjojsl123djioqwhfjnd"
	signingKey = "lasjdoiqjwdnkjsdhfmnasd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user serv.User) (int, error) {
	rand.Seed(time.Now().UnixNano())
	pass := ""
	for i := 0; i < 8; i++ {
		pass += string(passwordLetters[rand.Intn(82)])
	}
	/*
	---Отправка письма с паролем на почту 
	text := "Your password:\n    " + pass + "\nНикому не передавайте пароль!"
	if err := NewEmailService().SendEmail(user.Email, "Password", text); err != nil {
		return 0, err
	}*/
	log.Print("ПАРОЛЬ: ", pass)
	user.Password = geneartePasswordHash(pass)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repo.GetUser(email, geneartePasswordHash(password))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id_User,
	})

	return token.SignedString([]byte(signingKey))
}

func geneartePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
