package repository

import (
	serv "github.com/callmehorhe/backtest"
	"gorm.io/gorm"
)

type CafePostgres struct {
	db *gorm.DB
}

func NewCafePostgres(db *gorm.DB) *CafePostgres{
	return &CafePostgres{
		db: db,
	}
}

func (r *CafePostgres) GetCafeList() []serv.Cafe {
	var cafes []serv.Cafe
	r.db.Table("cafes").Find(&cafes)
	return cafes
}