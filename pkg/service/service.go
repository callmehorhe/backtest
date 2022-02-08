package service

import (
	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
)

type Authorization interface {
	CreateUser(user serv.User) (int, error)
	GenerateToken(email, password string) (string, error)
}

type EmailSendler interface {
	SendEmail(email, subject, text string) error
}

type Service struct {
	Authorization
	EmailSendler
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		EmailSendler: NewEmailService(),
	}
}
