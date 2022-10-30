package telegram

import (
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
	var status models.Status = models.Status(tmp[0])

	if order.Address == models.TakeawayOrder {
		switch status {
		case models.Accepted:
			b.cafeBot.Accept(callback.Message.Chat.ID, id, &order)
		case models.Ready:
			b.cafeBot.Ready(callback.Message.Chat.ID, id, &order)
		case models.Canceled:
			b.cafeBot.Cancel(callback.Message.Chat.ID, id, &order)
		case models.Sent:
			b.cafeBot.Send(callback.Message.Chat.ID, id, &order)
		default:
			logrus.Error("unknown state: %v", order.Status)
			return
		}
	} else {
		switch status {
		case models.Accepted:
			b.cafeBot.Accept(callback.Message.Chat.ID, id, &order)
			b.NewOrderForDrivers(order)
		case models.Sent:
			b.cafeBot.Send(callback.Message.Chat.ID, id, &order)
		case models.Canceled:
			b.cafeBot.Cancel(callback.Message.Chat.ID, id, &order)
		default:
			logrus.Error("unknown state: %v", order.Status)
			return
		}
	}

	b.repo.Orders.UpdateOrder(order)
}

func (b *CafeBot) Accept(chatId int64, id int, order *models.Order) {
	if order.Status != models.New {
		b.SendMessage(chatId, fmt.Sprintf("🛑Заказ №%d был подтвержден!", id))
		logrus.Error("order status not NEW")
		return
	}

	order.Status = models.Accepted
	b.SendMessage(chatId, fmt.Sprintf("✔️Заказ №%d подтвержден!", id))
}

func (b *CafeBot) Ready(chatId int64, id int, order *models.Order) {
	if order.Status != models.Accepted {
		b.SendMessage(chatId, fmt.Sprintf("🛑Для начала подтвердите заказ №%d!", id))
		logrus.Error("order status not ACCEPTED")
		return
	}
	order.Status = models.Ready
	b.SendMessage(chatId, fmt.Sprintf("✔️Заказ №%d готов!", id))
}

func (b *CafeBot) Send(chatId int64, id int, order *models.Order) {
	if order.Status != models.AcceptedByDriver && order.Status != models.Accepted {
		b.SendMessage(chatId, fmt.Sprintf("🛑Для начала подтвердите заказ №%d!", id))
		logrus.Errorf("order status not ACCEPTED or ACCEPTED_BY_DRIVER: %v", order.Status)
		return
	}
	order.Status = models.Sent
	b.SendMessage(chatId, fmt.Sprintf("✅Заказ №%d отправлен!", id))
}

func (b *CafeBot) Cancel(chatId int64, id int, order *models.Order) {
	if order.Status != models.New {
		b.SendMessage(chatId, fmt.Sprintf("🛑Невозможно отменить принятый заказ №%d!", id))
		logrus.Error("order status not ACCEPTED")
		return
	}
	order.Status = models.Canceled
	b.SendMessage(chatId, fmt.Sprintf("❌Заказ №%d отменен!", id))
}
