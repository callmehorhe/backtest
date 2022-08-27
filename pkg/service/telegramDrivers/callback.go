package telegramdrivers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (b *BotDrivers) NewOrder(order models.Order) {
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
	_, err := b.bot.Send(msg)
	if err != nil {
		logrus.Errorf("message cant be sent: %v", err)
		return
	}
}

func (b *BotDrivers) CallbackHandler(callback tgbotapi.CallbackQuery) {
	tmp := strings.Split(callback.Data, "f")

	id, _ := strconv.Atoi(tmp[1])
	switch tmp[0] {
	case "accept":
		b.Accept(callback, id)
	}
}

func (b *BotDrivers) Accept(callback tgbotapi.CallbackQuery, order_ID int) {
	order := b.repo.Orders.GetOrderByID(order_ID)
	if order.Driver != 0 {
		b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf("🛑Заказ №%d был подтвержден ранее водителем @%s!", order_ID, callback.From.UserName))
		return
	}
	order.Driver = int64(callback.From.ID)
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(-626247381, fmt.Sprintf("✔️Заказ №%d подтвержден водителем @%s!", order_ID, callback.From.UserName))
	b.SendFullOrder(order_ID, callback.From.ID)
}

func (b *BotDrivers) SendFullOrder(order_ID, driver_ID int) error {
	order := b.repo.Orders.GetOrderByID(order_ID)
	text := fmt.Sprintf("Заказ №%d\n%s\n", order.Order_ID, order.Cafe_Name)
	text += fmt.Sprintf("Адрес: %s\n", order.Address)
	text += fmt.Sprintf("📱Номер телефона: %s\n", order.Phone)
	for i := range order.Positions {
		text += fmt.Sprintf("%d: %s - %d шт.\n", i+1, order.Positions[i].Name, order.Positions[i].Count)
	}
	text += "Доставка: 100р\n"
	text += "Сервисный сбор: 20р\n"
	text += fmt.Sprintf("💸Итого: %dр.", order.Cost)
	msg := tgbotapi.NewMessage(int64(driver_ID), text)
	_, err := b.bot.Send(msg)
	if err != nil {
		logrus.Errorf("cant send message to tgDeliveryBot, %v", err)
		return err
	}

	return nil
}
