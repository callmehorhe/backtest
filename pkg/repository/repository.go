package repository

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	GetUserById(id int) (models.User, error)
	ConfirmUser(code string) error
}

type CafeList interface {
	GetCafeList() []models.Cafe
	GetMenuByCafeID(id int) []models.Menu
	GetCafeByID(id int) models.Cafe
	AddChatId(cafe_id int, chat_id int64)
	GetCategoriesByCafeID(id int) []string
	GetCafe(id int, password string) (models.Cafe, error)
	UpdateCafe(cafe models.Cafe) error
	CreatePos(menu models.Menu)
	UpdatePos(menu models.Menu)
	DeletePos(id []int)
	GetCafeChatId(id int) int64
}

type Orders interface {
	CreateOrder(order models.Order) int
	UpdateOrder(order models.Order) models.Order
	GetOrderByID(id int) models.Order
	GetOrdersByUser(id, cafe int) []models.Order
	GetCafeNameByID(id int) string
	GetOrdersCount(id int) int
}

type Drivers interface {
	IsNew(id int64) bool
	CreateDriver(driver models.Driver) error
	GetDriverById(id int64) (models.Driver, error)
}

type Repository struct {
	Authorization
	CafeList
	Orders
	Drivers
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		CafeList:      NewCafePostgres(db),
		Orders:        NewOrderPostgres(db),
		Drivers:       NewDriverPostgres(db),
	}
}
