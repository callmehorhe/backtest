package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	handler = ""
)

type Cafe struct {
	Name   string
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
		b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите имя организации: "))
		handler = "SignUp"
	case "info":
	}
}

func (b *Bot) SignUp(message *tgbotapi.Message) {
	cafe := Cafe{}
	cafe.Name = message.Text
	cafe.ChadID = message.Chat.ID
	b.bot.Send(tgbotapi.NewMessage(cafe.ChadID, "Вы успешно прошли регистрацию."))
}

