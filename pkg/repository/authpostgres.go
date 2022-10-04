package repository

import (
	"errors"

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
	err := r.db.Select("id_user, phone").Where("email=? AND password=? AND confirm=?", email, password, "").Take(&user).Error
	return user, err
}

func (r *AuthPostgres) GetUserById(id int) (models.User, error) {
	var user models.User
	err := r.db.Where("id_user=?", id).First(&user).Error
	return user, err

}

func (r *AuthPostgres) ConfirmUser(code string) error {
	var id int
	err := r.db.Select("id_user").Table("users").Where("confirm=?", code).Scan(&id).Error
	if err != nil || id == 0 {
		return errors.New("account doesn't confirm")
	}
	return r.db.Table("users").Where("id_user=?", id).Update("confirm", "").Error
}

func (r *AuthPostgres) ForgetPassword(email, phone, auth string) error {
	var user models.User
	if err := r.db.Table("users").Where("email=? AND phone=?", email, phone).Scan(&user).Error; err != nil {
		return err
	}
	user.Password = ""
	user.Confirm = auth
	if err := r.db.Table("users").Where("email=?", email).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) ResetPassword(auth, pass string) error {
	var user models.User
	if err := r.db.Table("users").Where("confirm=?").Scan(&user).Error; err != nil {
		return err
	}
	user.Confirm = ""
	user.Password = pass

	if err := r.db.Table("users").Where("id_user=?", user.Id_User).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}
