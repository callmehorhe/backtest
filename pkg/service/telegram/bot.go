package telegram

import (
	"github.com/callmehorhe/backtest/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	repo repository.CafeList
}

func NewBotService(repo repository.CafeList, bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot: bot,
		repo: repo,
	}
}

func (b *Bot) Start() error {	
	logrus.Printf("Bot %s activated", b.bot.Self.UserName)

	updates, err := b.initUpdateChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)
	return nil
}

func (b *Bot) initUpdateChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			b.HandleMessge(update.Message)
		}else if update.CallbackQuery != nil {
			b.CallbackHandler(*update.CallbackQuery)
		}
	}
}