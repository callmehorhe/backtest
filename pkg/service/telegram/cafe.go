package telegram

import (
	"fmt"
	"strconv"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Cafe struct {
	ID     int
	ChadID int64
}

func (b *CafeBot) SendMessage(chatID int64, text string) {
	b.bot.Send(tgbotapi.NewMessage(chatID, text))
}

func (b *Bot) StartCafe() error {
	logrus.Printf("Bot %s activated", b.cafeBot.bot.Self.UserName)

	updates, err := b.initUpdateChannel(b.cafeBot.bot)
	if err != nil {
		return err
	}
	b.handleUpdatesCafe(updates)
	return nil
}

func (b *Bot) handleUpdatesCafe(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			isNew := true
			for i := 0; i < len(cashers); i++ {
				if cashers[i].Id == update.Message.Chat.ID {
					b.HandleMessgeFromCafe(update.Message, &cashers[i])
					isNew = false
					break
				}
			}
			if isNew {
				casher := models.Cashers{
					Id: update.Message.Chat.ID,
				}
				b.HandleMessgeFromCafe(update.Message, &casher)
				cashers = append(cashers, casher)
			}
		} else if update.CallbackQuery != nil {
			b.CallbackHandlerForCafe(*update.CallbackQuery)
		}
	}
}

func (b *Bot) HandleMessgeFromCafe(message *tgbotapi.Message, casher *models.Cashers) {
	if message.IsCommand() {
		b.HandleCommandCafe(message, casher)
		return
	}
	if message.Text != "" {
		b.HandleTextCafe(message, casher)
	}
}

func (b *Bot) HandleCommandCafe(message *tgbotapi.Message, casher *models.Cashers) {
	switch message.Command() {
	case "start":
		b.cafeBot.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите ID:"))
		casher.Handler = "SignUpID"
	case "pass":
		b.cafeBot.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите ID:"))
	case "cancel":
		casher.Handler = ""
	case "info":
	}
}

func (b *Bot) HandleTextCafe(message *tgbotapi.Message, casher *models.Cashers) {
	switch casher.Handler {
	case "SignUpID":
		b.SignUpInCafeBot(message, casher)
	case "SignUpPass":
		b.SignUpPassBot(message, casher)
	}
}

func (b *Bot) SignUpInCafeBot(message *tgbotapi.Message, casher *models.Cashers) {
	cafeId, err := strconv.Atoi(message.Text)
	if err != nil {
		b.cafeBot.SendMessage(message.Chat.ID, "ID введен некорректно")
		return
	}
	cafe := b.repo.GetCafeByID(cafeId)
	if cafe.Id_Cafe == 0 {
		b.cafeBot.SendMessage(message.Chat.ID, "Введен несуществующий ID")
		return
	}

	//b.repo.AddChatId(cafe.Id_Cafe, message.Chat.ID)
	//b.cafeBot.SendMessage(message.Chat.ID, "Регистрация прошла успешно!")
	b.cafeBot.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите пароль:"))
	casher.Handler = "SignUpPass"
	casher.CafeID = cafeId
}

func (b *Bot) SignUpPassBot(message *tgbotapi.Message, casher *models.Cashers) {
	pass := message.Text
	casher.Handler = ""

	cafe, err := b.repo.CafeList.GetCafe(casher.CafeID, pass)
	if err != nil {
		b.cafeBot.SendMessage(casher.Id, "Введен неверный пароль!")
		return
	}

	cafe.Chat_ID = append(cafe.Chat_ID, message.Chat.ID)
	if err := b.repo.CafeList.UpdateCafe(cafe); err != nil {
		logrus.Errorf("[SignUpCafe] cant update cafe db: %v", err)
		b.cafeBot.SendMessage(casher.Id, "Невозможно привязать аккаунт! Обратитесь в поддержку!")
		return
	}

	b.cafeBot.SendMessage(casher.Id, fmt.Sprintf("Регистрация прошла успешно! Теперь вы будете получать информацию о заказах из кафе %v.", cafe.Name))
}

func (b *CafeBot) SendDriverInfo(driver models.Driver, cafeId int64, orderId int) {
	msg := fmt.Sprintf(`
		Заказ %d принят! Данные водителя:
		Имя: %s
		Машина: %s
		Номер телефона: %s
	`, orderId, driver.Name, driver.Car, driver.Phone)
	b.SendMessage(cafeId, msg)
}
