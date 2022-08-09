package main

import (
	"time"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	dataBase, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Errorf("cannot connect to db: %v", err)
	}
	db = dataBase
}

func ClearDB() {
	for {
		ClearUsers()

		time.Sleep(time.Hour * 24)
	}
}

func ClearUsers() {
	var users []models.User
	db.Table("users").Find(&users)
	for _, user := range users {
		if user.Confirm != "" {
			if err := db.Table("users").Delete(&user).Error; err != nil {
				logrus.Errorf("ClearUsers error: %v", err)
			}
		}
	}
}
