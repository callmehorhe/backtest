package repository

import (
	serv "github.com/callmehorhe/backtest"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user serv.User) (int, error)
	GetUser(username, password string) (serv.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *gorm.DB) *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}