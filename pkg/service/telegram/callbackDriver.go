package telegram

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (b *Bot) NewOrderForDrivers(order models.Order) {
	//Сообщение телеграмм бота
	text := fmt.Sprintf("Заказ №%d\n%s\n", order.Order_ID, order.Cafe_Name)
	text += fmt.Sprintf("Адрес: %s\n", order.Address)
	text += fmt.Sprintf("Сумма: %dр.", order.Cost)
	nKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Принять заказ", fmt.Sprintf("%sf%d", models.Accepted, order.Order_ID)),
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

func (b *Bot) CallbackHandler(callback tgbotapi.CallbackQuery) {
	tmp := strings.Split(callback.Data, "f")
	id, _ := strconv.Atoi(tmp[1])
	order := b.repo.Orders.GetOrderByID(id)

	driver, err := b.repo.Drivers.GetDriverById(int64(callback.From.ID))
	if err != nil {
		logrus.Error("cant find driver: %v", err)
		b.driverBot.SendMessage(-626247381, fmt.Sprintf("✔️Заказ №%d: ошибка подтверждения", order.Order_ID))
		return
	}
	logrus.Info(tmp[0])
	switch tmp[0] {
	case string(models.Accepted):
		b.driverBot.Accept(callback, id, &order, driver)
		cafechat := b.repo.CafeList.GetCafeChatId(order.Cafe_Id)
		b.cafeBot.SendDriverInfo(driver, cafechat, order.Order_ID)
	case string(models.Delivered):
		order.Status = models.Delivered
		b.driverBot.SendMessage(int64(callback.From.ID), "Заказ доставлен!")
	}
	b.repo.Orders.UpdateOrder(order)
}

func (b *DriverBot) Accept(callback tgbotapi.CallbackQuery, order_ID int, order *models.Order, driver models.Driver) {
	logrus.Info(order.Driver_Id)
	if order.Driver_Id != 0 {
		b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf("🛑Заказ №%d уже был подтвержден!", order_ID))
		return
	}
	order.Driver_Id = driver.Id
	if order.Status != models.Sent {
		order.Status = models.AcceptedByDriver
	}
	logrus.Warn("add driver to order: %+v", order)
	b.SendMessage(-626247381, fmt.Sprintf("✔️Заказ №%d подтвержден водителем %s %s!", order.Order_ID, callback.From.FirstName, callback.From.LastName))
	b.SendFullOrder(order, callback.From.ID)
	return
}

func (b *DriverBot) SendFullOrder(order *models.Order, driver_ID int) error {
	text := fmt.Sprintf("Заказ №%d\n%s\n", order.Order_ID, order.Cafe_Name)
	text += fmt.Sprintf("Адрес: %s\n", order.Address)
	text += fmt.Sprintf("📱Номер телефона: %s\n", order.Phone)
	order_list := make([]models.Position, 10)
	err := json.Unmarshal(order.Order_list, &order_list)
	if err != nil {
		logrus.Error("Cant unmarshal order_List: %v", err)
	}
	for i, pos := range order_list {
		text += fmt.Sprintf("%d: %s - %d шт.\n", i+1, pos.Name, pos.Count)
	}
	text += "Доставка: 100р\n"
	text += "Сервисный сбор: 20р\n"
	text += fmt.Sprintf("💸Итого: %dр.", order.Cost)
	nKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заказ доставлен", fmt.Sprintf("%sf%d", models.Delivered, order.Order_ID)),
		),
	)
	msg := tgbotapi.NewMessage(int64(driver_ID), text)
	msg.ReplyMarkup = nKeyboard
	_, err = b.bot.Send(msg)
	if err != nil {
		logrus.Errorf("cant send message to tgDeliveryBot, %v", err)
		return err
	}

	return nil
}
