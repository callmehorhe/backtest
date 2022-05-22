package repository

import (
	serv "github.com/callmehorhe/backtest"
	"gorm.io/gorm"
)

type CafePostgres struct {
	db *gorm.DB
}

func NewCafePostgres(db *gorm.DB) *CafePostgres {
	return &CafePostgres{
		db: db,
	}
}

func (r *CafePostgres) GetCafe(id int, password string) (serv.Cafe, error) {
	var cafe serv.Cafe
	err := r.db.Table("cafes").Where("id_cafes=? and password=?", id, password).Take(&cafe).Error
	return cafe, err
}

func (r *CafePostgres) GetCafeList() []serv.Cafe {
	var cafes []serv.Cafe
	r.db.Table("cafes").Find(&cafes)
	return cafes
}

func (r *CafePostgres) GetMenuByCafeID(id int) []serv.Menu {
	var positions []serv.Menu
	r.db.Table("menu").Where("id_cafe=?", id).Find(&positions)
	return positions
}

func (r *CafePostgres) GetCafeByID(id int) serv.Cafe {
	var cafe serv.Cafe
	r.db.Table("cafes").Where("id_cafe=?", id).First(&cafe)
	return cafe
}

func (r *CafePostgres) AddChatId(cafe_id int, chat_id int64) {
	r.db.Table("cafes").Where("id_cafe=?", cafe_id).Update("chat_id", chat_id)
}

func (r *CafePostgres) GetCategoriesByCafeID(id int) []string {
	var categories []string
	r.db.Table("menu").Select("category").Where("id_cafe=?", id).Group("category").Find(&categories)
	return categories
}
