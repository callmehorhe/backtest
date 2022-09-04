package telegramdrivers

import (
	"encoding/json"
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
		b.SendMessage(callback.Message.Chat.ID, fmt.Sprintf("🛑Заказ №%d уже был подтвержден!", order_ID))
		return
	}
	order.Driver = int64(callback.From.ID)
	b.repo.Orders.UpdateOrder(order)
	b.SendMessage(-626247381, fmt.Sprintf("✔️Заказ №%d подтвержден водителем %s %s!", order_ID, callback.From.FirstName, callback.From.LastName))
	b.SendFullOrder(order, callback.From.ID)
	cafechat := b.repo.CafeList.GetCafeChatId(order.Cafe_Id)
	b.SendDriverInfo(int64(callback.From.ID), cafechat, order_ID)
}

func (b *BotDrivers) SendFullOrder(order models.Order, driver_ID int) error {
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
	msg := tgbotapi.NewMessage(int64(driver_ID), text)
	_, err = b.bot.Send(msg)
	if err != nil {
		logrus.Errorf("cant send message to tgDeliveryBot, %v", err)
		return err
	}

	return nil
}

func (b *BotDrivers) SendDriverInfo(driverId, cafeId int64, orderId int) {
	driver, err := b.repo.Drivers.GetDriverById(driverId)
	if err != nil {
		logrus.Error("cant get driver: %v", err)
	}
	msg := fmt.Sprintf(`
		Заказ %d принят! Данные водителя:
		Имя: %s
		Машина: %s
		Номер телефона: %s
	`, orderId, driver.Name, driver.Car, driver.Phone)
	b.cafeBot.Send(tgbotapi.NewMessage(cafeId, msg))
}
