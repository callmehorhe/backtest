package telegram

import (
	"fmt"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (b *Bot) SendOrder(order models.Order) (models.Order, error) {
	if order.Phone == "" {
		user, _ := b.repo.GetUserById(order.User_ID)
		order.Phone = user.Phone
	}
	order.Status = models.New
	id := b.repo.Orders.CreateOrder(order)
	cafe := b.repo.GetCafeByID(order.Cafe_Id)
	order.Cafe_Name = cafe.Name

	// Сообщение телеграмм бота
	text := fmt.Sprintf("Заказ №%d\n%s\n", id, cafe.Name)

	// Высчитываем цену заказа и записываем ее сразу в order.Cost, а также получаем текст заказа для сообщения
	_, orderList := OrderCost(&order)
	text += orderList
	text += "Сервисный сбор: 20р\n"

	// Если адрес заказа не указан, считаем, что это вынос
	var nKeyboard tgbotapi.InlineKeyboardMarkup
	if order.Address != models.TakeawayOrder {
		text += fmt.Sprintf("📍Адрес: %s\n", order.Address)
		nKeyboard = DeliveryKeyboard(id)
	} else {
		nKeyboard = TakeawayKeyboard(id)
	}
	text += fmt.Sprintf("📱Номер телефона: %s\n", order.Phone)

	// Отправка сообщения в телеграмм всем кассирам
	for _, casher := range cafe.Chat_ID {
		msg := tgbotapi.NewMessage(casher, text)
		msg.ReplyMarkup = nKeyboard
		_, err := b.cafeBot.bot.Send(msg)
		if err != nil {
			logrus.Errorf("cant send message to tgDeliveryBot, %v", err)
			return models.Order{}, err
		}

	}
	order.Order_ID = id
	return b.repo.UpdateOrder(order), nil
}

func DeliveryKeyboard(id int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять заказ", fmt.Sprintf("%sf%d", models.Accepted, id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить заказ", fmt.Sprintf("%sf%d", models.Sent, id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить заказ", fmt.Sprintf("%sf%d", models.Canceled, id)),
		),
	)
}

func TakeawayKeyboard(id int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять заказ", fmt.Sprintf("%sf%d", models.Accepted, id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заказ готов", fmt.Sprintf("%sf%d", models.Ready, id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить заказ", fmt.Sprintf("%sf%d", models.Canceled, id)),
		),
	)
}

func OrderCost(order *models.Order) (int, string) {
	text := ""
	order.Cost = models.ServicePrice

	for i := range order.Positions {
		text += fmt.Sprintf("%d: %s - %d шт.\n", i+1, order.Positions[i].Name, order.Positions[i].Count)
		order.Cost += order.Positions[i].Count * order.Positions[i].Price
	}

	if order.Address == "" {
		text += "📌Навынос\n"
		order.Address = models.TakeawayOrder
	} else {
		text += "Доставка: 100р\n"
		order.Cost += models.DeliveryPrice
	}
	text += fmt.Sprintf("💸Итого: %dр.", order.Cost)

	return order.Cost, text
}
