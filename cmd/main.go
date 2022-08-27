package main

import (
	"os"

	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/handler"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/callmehorhe/backtest/pkg/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatal("fail init db")
	}

	tgBot, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		logrus.Fatal("cant launch bot")
	}
	tgDriverBot, err := tgbotapi.NewBotAPI("5527799879:AAFuVJfwGqKAcAmYFm7GKDNAjWsg_Dtpsko")
	if err != nil {
		logrus.Fatal("cant launch bot")
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos, tgBot, tgDriverBot)
	handlers := handler.NewHandler(services)
	go services.TGBot.Start()

	srv := serv.Server{}
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatal(err)
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
