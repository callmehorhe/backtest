package telegram

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	handler = ""
)

type Cafe struct {
	ID     int
	ChadID int64
}

func (b *Bot) HandleMessge(message *tgbotapi.Message) {
	if message.IsCommand() {
		b.HandleCommand(message)
		return
	}
	if message.Text != "" {
		b.HandleText(message)
	}
}

func (b *Bot) HandleText(message *tgbotapi.Message) {
	switch handler {
	case "SignUp":
		b.SignUp(message)
	}
}

func (b *Bot) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите ID:"))
		handler = "SignUp"
	case "cancel":
		handler = ""
	case "info":
	}
}

func (b *Bot) SignUp(message *tgbotapi.Message) {
	cafeId, err := strconv.Atoi(message.Text)
	if err != nil {
		b.SendMessage(message.Chat.ID, "ID введен некорректно")
		return
	}
	cafe := b.repo.GetCafeByID(cafeId)
	log.Print(cafe)
	if cafe.Id_Cafe == 0 {
		b.SendMessage(message.Chat.ID, "Введен несуществующий ID")
		return
	}

	if cafe.Chat_ID != 0 {
		b.SendMessage(message.Chat.ID, "ID уже занят.")
		b.SendMessage(cafe.Chat_ID, "Попытка повторной привязки вашего ID к другому устройству!")
		return
	}

	b.repo.AddChatId(cafe.Id_Cafe, message.Chat.ID)
	b.SendMessage(message.Chat.ID, "Регистрация прошла успешно!")
	handler = ""
}

func (b *Bot) SendMessage(chatID int64, text string) {
	b.bot.Send(tgbotapi.NewMessage(chatID, text))
}
