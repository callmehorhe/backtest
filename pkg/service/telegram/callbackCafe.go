package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (b *Bot) CallbackHandlerForCafe(callback tgbotapi.CallbackQuery) {
	tmp := strings.Split(callback.Data, "f")
	id, _ := strconv.Atoi(tmp[1])
	order := b.repo.Orders.GetOrderByID(id)

	switch tmp[0] {
	case "accept":
		updatedOrder, err := b.cafeBot.Accept(callback.Message.Chat.ID, id, order)
		if err != nil {
			logrus.Error(err)
		}

		b.repo.Orders.UpdateOrder(updatedOrder)

		if updatedOrder.Address != "Навынос" {
			b.NewOrderForDrivers(order)
		}
	case "send":
		b.Send(callback.Message.Chat.ID, id)
	case "cancel":
		b.Cancel(callback.Message.Chat.ID, id)
	}
}

func (b *CafeBot) Accept(chat_ID int64, id int, order models.Order) (models.Order, error) {
	if order.Status == "acceptet" {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был подтвержден ранее!", id))
		return models.Order{}, errors.New("order was accepted")
	} else if order.Status == "canceled" {
		b.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был отменен!", id))
		return models.Order{}, errors.New("order was canceled")
	}
	
	order.Status = "accepted"

	b.SendMessage(chat_ID, fmt.Sprintf("✔️Заказ №%d подтвержден!", id))
	return order, nil
}

func (b *Bot) Send(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if order.Status != "accepted" {
		b.cafeBot.SendMessage(chat_ID, fmt.Sprintf("🛑Для начала подтвердите заказ №%d!", id))
		return
	} else if order.Status == "canceled" {
		b.cafeBot.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был отменен!", id))
		return
	}
	order.Status = "sent"
	b.repo.Orders.UpdateOrder(order)
	b.cafeBot.SendMessage(chat_ID, fmt.Sprintf("✅Заказ №%d отправлен!", id))
}

func (b *Bot) Cancel(chat_ID int64, id int) {
	order := b.repo.Orders.GetOrderByID(id)
	if order.Status == "accepted" || order.Status == "sent" {
		b.cafeBot.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d уже был подтвержден! Нельзя отменить подтвержденный заказ!", id))
		return
	} else if order.Status == "canceled" {
		b.cafeBot.SendMessage(chat_ID, fmt.Sprintf("🛑Заказ №%d был отменен ранее!", id))
		return
	}
	order.Status = "canceled"
	b.repo.Orders.UpdateOrder(order)
	b.cafeBot.SendMessage(chat_ID, fmt.Sprintf("❌Заказ №%d отменен!", id))
}
