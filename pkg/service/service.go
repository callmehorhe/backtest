package service

import (
	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/callmehorhe/backtest/pkg/service/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	GetMenuByCafeID(id int) []serv.Menu
	GetCafeByID(id int) serv.Cafe
	GetCategoriesByCafeID(id int) []string
}

type TGBot interface {
	Start() error
	HandleMessge(message *tgbotapi.Message)
	SendOrder(order serv.Order)
}

type Service struct {
	Authorization
	EmailSendler
	CafeList
	TGBot
}



func NewService(repos *repository.Repository, bot *tgbotapi.BotAPI) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		EmailSendler: NewEmailService(),
		CafeList: NewCafeService(repos.CafeList),
		TGBot: telegram.NewBotService(*repos, bot),
	}
}
