package service

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/callmehorhe/backtest/pkg/service/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	GetUser(email, password string) (models.User, error)
	GetUserByID(id int) (models.User, error)
}

type EmailSendler interface {
	SendEmail(email, subject, text string) error
}

type CafeList interface {
	GetCafeList() []models.Cafe
	GetMenuByCafeID(id int) []models.Menu
	GetCafeByID(id int) models.Cafe
	GetCategoriesByCafeID(id int) []string
	GetCafe(id int, password string) (models.Cafe, error)
	CafeGenerateToken(id int, password string) (string, error)
	UpdateCafe(cafe models.Cafe) error
	UpdateMenu(menu []models.Menu)
}

type TGBot interface {
	Start() error
	HandleMessge(message *tgbotapi.Message)
	SendOrder(order models.Order)
}

type Order interface {
	GetOrdersByUser(id int) []models.Order
}

type Service struct {
	Authorization
	EmailSendler
	CafeList
	TGBot
	Order
}

func NewService(repos *repository.Repository, bot *tgbotapi.BotAPI) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		EmailSendler:  NewEmailService(),
		CafeList:      NewCafeService(repos.CafeList),
		TGBot:         telegram.NewBotService(*repos, bot),
		Order:         NewOrderService(repos.Orders),
	}
}
