package telegramdrivers

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/callmehorhe/backtest/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type BotDrivers struct {
	bot     *tgbotapi.BotAPI
	repo    repository.Repository
	cafeBot *tgbotapi.BotAPI
}

var drivers []models.Driver

func NewBotService(repo repository.Repository, bot *tgbotapi.BotAPI, cafeBot *tgbotapi.BotAPI) *BotDrivers {
	return &BotDrivers{
		bot:  bot,
		repo: repo,
		cafeBot: cafeBot,
	}
}

func (b *BotDrivers) Start() error {
	logrus.Printf("Bot %s activated", b.bot.Self.UserName)
	drivers = make([]models.Driver, 0)
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
		logrus.Print("drivers: %+v", drivers)
		if update.Message != nil { // ignore any non-Message Updates
			isNew := true
			for i := 0; i < len(drivers); i++ {
				if drivers[i].Id == update.Message.Chat.ID {
					b.HandleMessge(update.Message, &drivers[i])
					isNew = false
					break
				}
			}
			if isNew {
				driver := models.Driver{
					Id: update.Message.Chat.ID,
				}
				b.HandleMessge(update.Message, &driver)
				drivers = append(drivers, driver)
			}
		} else if update.CallbackQuery != nil {
			b.CallbackHandler(*update.CallbackQuery)
		}
	}
}
