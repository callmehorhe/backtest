package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	passwordLetters = "0123456789abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt            = "opie435qjojsl123djioqwhfjnd"
	signingKey      = "lasjdoiqjwdnkjsdhfmnasd"
	tokenTTL        = 12 * time.Hour
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

func (s *AuthService) CreateUser(user models.User) (int, error) {
	if len(user.Password) < 8 {
		return 0, errors.New("password too short")
	}
	rand.Seed(time.Now().UnixNano())
	auth := ""
	for i := 0; i < 20; i++ {
		auth += string(passwordLetters[rand.Intn(82)])
	}
	user.Confirm = auth
	user.Password = geneartePasswordHash(user.Password)
	user.Email = strings.ToLower(user.Email)
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, err
	}
	text := fmt.Sprintf("Ссылка для подтверждения:\nhttp://dhouse-it.ru/confirm/%s", user.Confirm)
	if err := NewEmailService().SendEmail(user.Email, "Password", text); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *AuthService) GetUserByID(id int) (models.User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil || user.Confirm != "" {
		return models.User{}, err
	}
	return user, nil
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

func (s *AuthService) GetUser(email, password string) (models.User, error) {
	email = strings.ToLower(email)
	return s.repo.GetUser(email, geneartePasswordHash(password))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func geneartePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) ConfirmUser(code string) error {
	return s.repo.ConfirmUser(code)
}
