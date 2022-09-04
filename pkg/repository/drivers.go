package repository

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"gorm.io/gorm"
)

type DriversPostgres struct {
	db *gorm.DB
}

func NewDriverPostgres(db *gorm.DB) *DriversPostgres {
	return &DriversPostgres{
		db: db,
	}
}

func (r *DriversPostgres) IsNew(id int64) bool {
	err := r.db.Select("id").Table("drivers").Where("id = ?", id).Error
	if err != nil {
		return false
	}
	return true
}
func (r *DriversPostgres) CreateDriver(driver models.Driver) error {
	return r.db.Create(&driver).Table("drivers").Error
}

func (r *DriversPostgres) GetDriverById(id int64) (models.Driver, error) {
	driver := models.Driver{}
	err := r.db.Table("drivers").Where("id = ?", id).Take(&driver).Error
	if err != nil {
		return models.Driver{}, err
	}
	return driver, nil
}