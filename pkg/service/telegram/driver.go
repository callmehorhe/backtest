package telegram

import (
	"fmt"

	"github.com/callmehorhe/backtest/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (b *DriverBot) SendMessage(chatID int64, text string) {
	b.bot.Send(tgbotapi.NewMessage(chatID, text))
}

func (b *Bot) StartDrivers() error {
	logrus.Printf("Bot %s activated", b.driverBot.bot.Self.UserName)

	updates, err := b.initUpdateChannel(b.driverBot.bot)
	if err != nil {
		return err
	}
	b.handleUpdatesDrivers(updates)
	return nil
}

func (b *Bot) handleUpdatesDrivers(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			isNew := true
			for i := 0; i < len(drivers); i++ {
				if drivers[i].Id == update.Message.Chat.ID {
					b.HandleMessgeFromDriver(update.Message, &drivers[i])
					isNew = false
					break
				}
			}
			if isNew {
				driver := models.Driver{
					Id: update.Message.Chat.ID,
				}
				b.HandleMessgeFromDriver(update.Message, &driver)
				drivers = append(drivers, driver)
			}
		} else if update.CallbackQuery != nil {
			b.CallbackHandler(*update.CallbackQuery)
		}
	}
}

func (b *Bot) HandleMessgeFromDriver(message *tgbotapi.Message, driver *models.Driver) {
	if message.IsCommand() {
		b.HandleCommandDriver(message, driver)
		return
	}
	
	if message.Text != "" {
		b.HandleTextDriver(message, driver)
	}
}

func (b *Bot) HandleCommandDriver(message *tgbotapi.Message, driver *models.Driver) {
	switch message.Command() {
	case "start":
		if b.repo.Drivers.IsNew(message.Chat.ID) {
			b.driverBot.SendMessage(message.Chat.ID, "Введите ФИО:")
			driver.Handler = "SignUp"
			return
		}
		b.driverBot.SendMessage(message.Chat.ID, "Вы уже зарегистрированы!")
	case "cancel":
		driver.Handler = ""
	case "info":
	}
}

func (b *Bot) HandleTextDriver(message *tgbotapi.Message, driver *models.Driver) {
	switch driver.Handler {
	case "SignUp":
		b.driverBot.Name(message, driver)
	case "CarModel":
		b.driverBot.CarModer(message, driver)
	case "CarNumber":
		b.driverBot.CarNumber(message, driver)
	case "Phone":
		driver := b.driverBot.Phone(message, driver)
		err := b.SignUpDriver(driver)
		msg := "Не удалось создать пользователя!"
		if err == nil {
			msg = fmt.Sprintf(`Вы прошли регистрацию! Ваши данные:
			Имя: %s
			Машина: %s
			Номер телефона: %s`,
				driver.Name, driver.Car, driver.Phone)
		} else {
			logrus.Error("cant create new driver: %v", err)
		}
		b.driverBot.SendMessage(message.Chat.ID, msg)
	}
}

func (b *Bot) SignUpDriver(driver models.Driver) error {
	err := b.repo.Drivers.CreateDriver(driver)
	if err != nil {
		b.driverBot.SendMessage(driver.Id, "Не удалось добавить водителя!")
		logrus.Error("cant add driver: %v", err)
		return err
	}
	return nil
}

func (b *DriverBot) Name(message *tgbotapi.Message, driver *models.Driver) {
	driver.Name = message.Text
	b.SendMessage(message.Chat.ID, "Введите марку и цвет машины: ")
	driver.Handler = "CarModel"
}

func (b *DriverBot) CarModer(message *tgbotapi.Message, driver *models.Driver) {
	driver.Car = message.Text
	b.SendMessage(message.Chat.ID, "Введите номер машины: ")
	driver.Handler = "CarNumber"
}

func (b *DriverBot) CarNumber(message *tgbotapi.Message, driver *models.Driver) {
	driver.Car += " " + message.Text
	b.SendMessage(message.Chat.ID, "Введите номер телефона в формате +7хххххххххх: ")
	driver.Handler = "Phone"
}

func (b *DriverBot) Phone(message *tgbotapi.Message, driver *models.Driver) models.Driver {
	driver.Phone = message.Text
	driver.Handler = ""
	return *driver
}
