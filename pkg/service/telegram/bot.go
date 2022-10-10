package telegram

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/callmehorhe/backtest/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var drivers []models.Driver
var cashers []models.Cashers

type CafeBot struct {
	bot *tgbotapi.BotAPI
}

type DriverBot struct {
	bot *tgbotapi.BotAPI
}

type Bot struct {
	cafeBot   CafeBot
	driverBot DriverBot
	repo      *repository.Repository
}

func NewBotService(repo repository.Repository, cafeBot, driverBot *tgbotapi.BotAPI) *Bot {
	cafe := CafeBot{
		bot: cafeBot,
	}
	driver := DriverBot{
		bot: driverBot,
	}
	return &Bot{
		cafeBot:   cafe,
		driverBot: driver,
		repo:      &repo,
	}
}

func (b *Bot) Start() error {
	go b.StartCafe()
	go b.StartDrivers()
	return nil
}

func (b *Bot) initUpdateChannel(bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return bot.GetUpdatesChan(u)
}
