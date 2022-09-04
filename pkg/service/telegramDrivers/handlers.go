package telegramdrivers

import (
	"fmt"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (b *BotDrivers) HandleMessge(message *tgbotapi.Message, driver *models.Driver) {
	if message.IsCommand() {
		b.HandleCommand(message, driver)
		return
	}
	if message.Text != "" {
		b.HandleText(message, driver)
	}
}

func (b *BotDrivers) HandleText(message *tgbotapi.Message, driver *models.Driver) {
	switch driver.Handler {
	case "SignUp":
		b.Name(message, driver)
	case "CarModel":
		b.CarModer(message, driver)
	case "CarNumber":
		b.CarNumber(message, driver)
	case "Phone":
		b.Phone(message, driver)
	}
}

func (b *BotDrivers) HandleCommand(message *tgbotapi.Message, driver *models.Driver) {
	switch message.Command() {
	case "start":
		if b.repo.Drivers.IsNew(message.Chat.ID) {
			b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Введите ФИО:"))
			driver.Handler = "SignUp"
			return
		}
		b.SendMessage(message.Chat.ID, "Вы уже зарегистрированы!")
	case "cancel":
		driver.Handler = ""
	case "info":
	}
}

func (b *BotDrivers) Name(message *tgbotapi.Message, driver *models.Driver) {
	driver.Name = message.Text
	b.SendMessage(message.Chat.ID, "Введите марку и цвет машины: ")
	driver.Handler = "CarModel"
}

func (b *BotDrivers) CarModer(message *tgbotapi.Message, driver *models.Driver) {
	driver.Car = message.Text
	b.SendMessage(message.Chat.ID, "Введите номер машины: ")
	driver.Handler = "CarNumber"
}

func (b *BotDrivers) CarNumber(message *tgbotapi.Message, driver *models.Driver) {
	driver.Car += " " + message.Text
	b.SendMessage(message.Chat.ID, "Введите номер телефона в формате +7хххххххххх: ")
	driver.Handler = "Phone"
}

func (b *BotDrivers) Phone(message *tgbotapi.Message, driver *models.Driver) {
	driver.Phone = message.Text
	msg := fmt.Sprintf("Вы прошли регистрацию! Ваши данные:\nИмя: %s\nМашина: %s\nНомер телефона: %s", driver.Name, driver.Car, driver.Phone)
	b.SendMessage(message.Chat.ID, msg)
	err := b.repo.Drivers.CreateDriver(*driver)
	if err != nil {
		b.SendMessage(message.Chat.ID, "Не удалось добавить водителя!")
		logrus.Error("cant add driver: %v", err)
	}
	driver.Handler = ""
}

func (b *BotDrivers) SendMessage(chatID int64, text string) {
	b.bot.Send(tgbotapi.NewMessage(chatID, text))
}
