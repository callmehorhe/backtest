package repository

import (
	"github.com/callmehorhe/backtest/pkg/models"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	row := r.db.Create(&user)
	if err := row.Error; err != nil {
		return 0, err
	}
	row.Select("id_user").Scan(&id)
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (models.User, error) {
	var user models.User
	user.Email = email
	user.Password = password
	err := r.db.Select("id_user").Where("email=? AND password=?", email, password).Take(&user).Error
	return user, err
}

func (r *AuthPostgres) GetUserById(id int) (models.User, error) {
	var user models.User
	err := r.db.Where("id_user=?", id).First(&user).Error
	return user, err

}
