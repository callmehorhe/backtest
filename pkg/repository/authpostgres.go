package repository

import (
	serv "github.com/callmehorhe/backtest"
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

func (r *AuthPostgres) CreateUser(user serv.User) (int, error) {
	var id int
	row := r.db.Create(&user)
	if err := row.Error; err != nil {
		return 0, err
	}
	row.Select("id_user").Scan(&id)
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) serv.User {
	var user serv.User
	user.Email = email
	user.Password = password
	r.db.Select("id_user").Where("email=? AND password=?", email, password).Take(&user)
	return user
}

func (r *AuthPostgres) GetUserById(id int) (serv.User, error) {
	var user serv.User
	err := r.db.Where("id_user=?", id).First(&user).Error
	return user, err

}
