package telegramdrivers

import (
	"github.com/callmehorhe/backtest/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type BotDrivers struct {
	bot *tgbotapi.BotAPI
	repo repository.Repository
}

func NewBotService(repo repository.Repository, bot *tgbotapi.BotAPI) *BotDrivers {
	return &BotDrivers{
		bot: bot,
		repo: repo,
	}
}

func (b *BotDrivers) Start() error {	
	logrus.Printf("Bot %s activated", b.bot.Self.UserName)

	updates, err := b.initUpdateChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)
	return nil
}

func (b *BotDrivers) initUpdateChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *BotDrivers) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			b.HandleMessge(update.Message)
		}else if update.CallbackQuery != nil {
			b.CallbackHandler(*update.CallbackQuery)
		}
	}
}