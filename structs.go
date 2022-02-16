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
	Chat_ID int64  `json:"-"`
}

type Menu struct {
	Id_Menu int    `gorn:"primary_key" json:"id"`
	Name    string `json:"productName"`
	Image   string `json:"img"`
	Price   int    `json:"cost"`
}

type CafeAndMenu struct {
	CafeName string `json:"name"`
	Menu     []Menu `json:"menu"`
}

type Position struct {
	Name  string `json:"name"`
	Price int    `json:"cost"`
	Count int    `json:"qty"`
	Sum   int    `json:"sum"`
}

type Order struct {
	CafeId  int        `json:"id"`
	Dishes  []Position `json:"menu"`
	Address string     `json:"address"`
}
