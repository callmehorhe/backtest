package repository

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/sirupsen/logrus"
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

func (r *CafePostgres) GetCafe(id int, password string) (models.Cafe, error) {
	var cafe models.Cafe
	err := r.db.Table("cafes").Where("id_cafe=? and password=?", id, password).Take(&cafe).Error
	return cafe, err
}

func (r *CafePostgres) GetCafeList() []models.Cafe {
	var cafes []models.Cafe
	r.db.Table("cafes").Order("enable desc, queue").Find(&cafes)
	logrus.Infof("cafes: %+v", cafes)
	return cafes
}

func (r *CafePostgres) GetMenuByCafeID(id int) []models.Menu {
	var positions []models.Menu
	r.db.Table("menu").Where("id_cafe=?", id).Find(&positions)
	return positions
}

func (r *CafePostgres) GetCafeByID(id int) models.Cafe {
	var cafe models.Cafe
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

func (r *CafePostgres) UpdateCafe(cafe models.Cafe) error {
	return r.db.Table("cafes").Where("id_cafe=?", cafe.Id_Cafe).Updates(&cafe).Error
}

func (r *CafePostgres) CreatePos(menu models.Menu) {
	m := models.Menu{
		Id_Cafe:     menu.Id_Cafe,
		Name:        menu.Name,
		Image:       menu.Image,
		Price:       menu.Price,
		Category:    menu.Category,
		Description: menu.Description,
		Weight:      menu.Weight,
		Avaible:     menu.Avaible,
	}
	r.db.Table("menu").Create(&m)
}

func (r *CafePostgres) UpdatePos(menu models.Menu) {
	r.db.Table("menu").Where("id_menu=?", menu.Id_Menu).Updates(&menu)
}

func (r *CafePostgres) DeletePos(id []int) {
	r.db.Table("menu").Delete(&models.Menu{}, id)
}

func (r *CafePostgres) GetCafeChatId(id int) int64 {
	var chat_id int64
	r.db.Table("cafes").Where("id_cafe = ?", id).Select("chat_id").Take(&chat_id)
	return chat_id
}
