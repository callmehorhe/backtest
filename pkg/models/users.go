package models

type User struct {
	Id_User  int    `json:"-"        gorm:"type:serial;primary_key"`
	Name     string `json:"name"     gorm:"type:varchar(255);not null"`
	Email    string `json:"email"    gorm:"type:varchar(255);unique;not null"`
	Phone    string `json:"phone"    gorm:"type:varchar(255);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}