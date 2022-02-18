package telegram

import (
	"fmt"
	"strconv"
	"strings"

	serv "github.com/callmehorhe/backtest"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) SendOrder(order serv.Order) {
	id := b.repo.Orders.CreateOrder(order)
	cafe := b.repo.GetCafeByID(order.Cafe_Id)
	text := fmt.Sprintf("Заказ №%d\n*%s*\n", id, cafe.Name)
	if order.Address != "" {
		text += fmt.Sprintf("📍Адрес: %s\n", order.Address)
	} else {
		text += "📌Навынос\n"
	}
	text += fmt.Sprintf("📱Номер телефона: %s\n", order.Phone)
	sum := 0
	for i := range order.Positions {
		text += fmt.Sprintf("%d: %s - %d шт.\n", i+1, order.Positions[i].Name, order.Positions[i].Count)
		sum += order.Positions[i].Count * order.Positions[i].Price
	}
	text += fmt.Sprintf("💸Итого: %dр.", sum)
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
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был подтвержден ранее!", id))
		return
	} else if order.Status_canceled {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был отменен!", id))
		return
	}
	order.Status_accepted = true
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(chat_ID, fmt.Sprintf("✔️Заказ №%d подтвержден!", id))
}

func (b *Bot) Send(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if !order.Status_accepted {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Для начала подтвердите заказ №%d!", id))
		return
	} else if order.Status_canceled {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был отменен!", id))
		return
	}
	order.Status_sent = true
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(chat_ID, fmt.Sprintf("✅Заказ №%d отправлен!", id))
}

func (b *Bot) Cancel(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if order.Status_accepted {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d уже был подтвержден! Нельзя отменить подтвержденный заказ!", id))
		return
	} else if order.Status_canceled {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был отменен ранее!", id))
		return
	}
	order.Status_canceled = true
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(chat_ID, fmt.Sprintf("❌Заказ №%d отменен!", id))
}
