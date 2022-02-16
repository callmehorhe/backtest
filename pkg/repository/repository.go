package repository

import (
	serv "github.com/callmehorhe/backtest"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user serv.User) (int, error)
	GetUser(username, password string) (serv.User, error)
	GetUserById(id int) (serv.User, error)
}

type CafeList interface {
	GetCafeList() []serv.Cafe
	GetMenuByCafeID(id int) []serv.Menu
	GetCafeByID(id int) serv.Cafe
	AddChatId(cafe_id int, chat_id int64)
}

type Repository struct {
	Authorization
	CafeList
}

func NewRepository(db *gorm.DB) *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
		CafeList: NewCafePostgres(db),
	}
}