package service

import (
	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
)

type Authorization interface {
	CreateUser(user serv.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	GetUserByID(id int) (serv.User, error)
}

type EmailSendler interface {
	SendEmail(email, subject, text string) error
}

type CafeList interface {
	GetCafeList() []serv.Cafe
}

type Service struct {
	Authorization
	EmailSendler
	CafeList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		EmailSendler: NewEmailService(),
		CafeList: NewCafeService(repos.CafeList),
	}
}
