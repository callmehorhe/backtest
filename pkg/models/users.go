package models

type User struct {
	Id_User  int    `json:"-"                gorm:"type:serial;primary_key"`
	Name     string `json:"name"             gorm:"type:varchar(255);not null"`
	Email    string `json:"email"            gorm:"type:varchar(255);unique;not null"`
	Phone    string `json:"phone"            gorm:"type:varchar(255);unique;not null"`
	Confirm  string `json:"authorization"    gorm:"type:varchar(255)"`
	Password string `json:"password"         gorm:"type:varchar(255);not null"`
}

type Driver struct {
	Id      int64  `json:"-"     gorm:"type:serial"`
	Name    string `json:"name"  gorm:"type:varchar(255);not null"`
	Car     string `json:"car"   gorm:"type:varchar(255);not null"`
	Phone   string `json:"phone" gorm:"type:varchar(255);not null"`
	Handler string `json:"-"     gorm:"-"`
}
