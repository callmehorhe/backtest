package telegram

import (
	"fmt"
	"strconv"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

var (
	handler = ""
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
			b.HandleMessgeFromCafe(update.Message)
		} else if update.CallbackQuery != nil {
			b.CallbackHandlerForCafe(*update.CallbackQuery)
		}
	}
}

func (b *Bot) HandleMessgeFromCafe(message *tgbotapi.Message) {
	if message.IsCommand() {
		b.HandleCommandCafe(message)
		return
	}
	if message.Text != "" {
		b.HandleTextCafe(message)
	}
}

func (b *Bot) HandleCommandCafe(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.cafeBot.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите ID:"))
		handler = "SignUp"
	case "cancel":
		handler = ""
	case "info":
	}
}

func (b *Bot) HandleTextCafe(message *tgbotapi.Message) {
	switch handler {
	case "SignUp":
		b.SignUpInCafeBot(message)
	}
}

func (b *Bot) SignUpInCafeBot(message *tgbotapi.Message) {
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

	if cafe.Chat_ID != 0 {
		b.cafeBot.SendMessage(message.Chat.ID, "ID уже занят.")
		b.cafeBot.SendMessage(cafe.Chat_ID, "Попытка повторной привязки вашего ID к другому устройству!")
		return
	}

	b.repo.AddChatId(cafe.Id_Cafe, message.Chat.ID)
	b.cafeBot.SendMessage(message.Chat.ID, "Регистрация прошла успешно!")
	handler = ""
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
