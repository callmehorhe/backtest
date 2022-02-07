package repository

import (
	"github.com/callmehorhe/backtest"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres{
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user serv.User)(int, error){
	var id int
	row := r.db.Create(&user)
	if err := row.Error; err != nil{
		return 0, err
	}
	row.Select("id_user").Scan(&id)
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (serv.User, error){
	var user serv.User
	err := r.db.Where("email = ? AND password = ?", email, password).Take(&user).Error
	return user, err
}