package serv

type User struct {
	Id_User  int    `gorm:"primary_key" json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Cafe struct {
	Id_Cafe int    `gorm:"primary_key" json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Image   string `json:"img"`
	Address string `json:"address"`
	Id_Menu int    `json:"menu_id"`
}

type Position struct {
	Name  string `json:"productName"`
	Price int    `json:"price"`
	Count int    `json:"qty"`
	Sum   int    `json:"sum"`
}
