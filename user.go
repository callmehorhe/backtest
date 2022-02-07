package serv

type User struct {
	Id_User  int    `gorm:"primary_key" json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
