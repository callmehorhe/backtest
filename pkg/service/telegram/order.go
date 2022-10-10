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

	id := b.repo.Orders.CreateOrder(order)
	cafe := b.repo.GetCafeByID(order.Cafe_Id)
	order.Cafe_Name = cafe.Name
	//Сообщение телеграмм бота
	text := fmt.Sprintf("Заказ №%d\n%s\n", id, cafe.Name)
	if order.Address != "" {
		text += fmt.Sprintf("📍Адрес: %s\n", order.Address)
		order.Cost += 100 //цена доставки
	} else {
		text += "📌Навынос\n"
		order.Address = "Навынос"
	}
	text += fmt.Sprintf("📱Номер телефона: %s\n", order.Phone)
	order.Cost = 20
	for i := range order.Positions {
		text += fmt.Sprintf("%d: %s - %d шт.\n", i+1, order.Positions[i].Name, order.Positions[i].Count)
		order.Cost += order.Positions[i].Count * order.Positions[i].Price
	}
	if order.Address != "Навынос" {
		text += "Доставка: 100р\n"
		order.Cost += 100
	}
	text += "Сервисный сбор: 20р\n"
	text += fmt.Sprintf("💸Итого: %dр.", order.Cost)
	nKeyboard := tgbotapi.NewInlineKeyboardMarkup(
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

func (b *Bot) NewOrderForDrivers(order models.Order) {
	//Сообщение телеграмм бота
	text := fmt.Sprintf("Заказ №%d\n%s\n", order.Order_ID, order.Cafe_Name)
	text += fmt.Sprintf("Адрес: %s\n", order.Address)
	text += fmt.Sprintf("Сумма: %dр.", order.Cost)
	nKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять заказ", fmt.Sprintf("acceptf%d", order.Order_ID)),
		),
	)
	msg := tgbotapi.NewMessage(-626247381, text)
	msg.ReplyMarkup = nKeyboard
	_, err := b.driverBot.bot.Send(msg)
	if err != nil {
		logrus.Errorf("message cant be sent: %v", err)
		return
	}
}
