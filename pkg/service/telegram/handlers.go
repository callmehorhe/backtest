package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	serv "github.com/callmehorhe/backtest"
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

func (b *Bot) SendOrder(order serv.Order) {
	id := b.repo.Orders.CreateOrder(order)
	cafe := b.repo.GetCafeByID(order.Cafe_Id)
	text := fmt.Sprintf("Заказ № %d\n%s\n", id, cafe.Name)
	sum := 0
	for i := range order.Positions {
		text += fmt.Sprintf("%d: %s - %d шт.\n", i+1, order.Positions[i].Name, order.Positions[i].Count)
		sum += order.Positions[i].Sum
	}
	text += fmt.Sprintf("Итого: %dр.", sum)
	nKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять заказ", fmt.Sprintf("acceptf%d", id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить заказ", fmt.Sprintf("sendf%d", id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить заказ", fmt.Sprintf("cancelf%d", id)),
		),
	)
	msg := tgbotapi.NewMessage(cafe.Chat_ID, text)
	msg.ReplyMarkup = nKeyboard
	b.bot.Send(msg)
}

func (b *Bot) CallbackHandler(callback tgbotapi.CallbackQuery) {
	tmp := strings.Split(callback.Data, "f")
	id, _ := strconv.Atoi(tmp[1])
	switch tmp[0] {
	case "accept":
		b.Accept(callback.Message.Chat.ID, id)
	case "send":
		b.Send(callback.Message.Chat.ID, id)
	case "cancel":
		b.Cancel(callback.Message.Chat.ID, id)
	}
}

func (b *Bot) Accept(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if order.Status_accepted {
		b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d был подтвержден ранее!", id))
		return
	} else if order.Status_canceled {
		b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d был отменен!", id))
		return
	}
	order.Status_accepted = true
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d подтвержден!", id))
}

func (b *Bot) Send(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if !order.Status_accepted {
		b.SendMessage(chat_ID, fmt.Sprintf("Для начала подтвердите заказ №%d!", id))
		return
	} else if order.Status_canceled {
		b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d был отменен!", id))
		return
	}
	order.Status_sent = true
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d отправлен!", id))
}

func (b *Bot) Cancel(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if order.Status_accepted {
		b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d уже был подтвержден! Нельзя отменить подтвержденный заказ!", id))
		return
	} else if order.Status_canceled {
		b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d был отменен ранее!", id))
		return
	}
	order.Status_canceled = true
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(chat_ID, fmt.Sprintf("Заказ №%d отменен!", id))
}
