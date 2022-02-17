package serv

import "gorm.io/datatypes"

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
	Id_Menu  int    `gorm:"primary_key" json:"id"`
	Name     string `json:"productName"`
	Image    string `json:"img"`
	Price    int    `json:"cost"`
	Category string `json:"category"`
}

type Category struct {
	Category_Name string `json:"name"`
	Menu_List     []Menu `json:"cart"`
}

type CafeAndMenu struct {
	Cafe_Name  string     `json:"name"`
	Categories []Category `json:"category"`
}

type Position struct {
	ID    int    `json:"productId"`
	Name  string `json:"productName"`
	Price int    `json:"cost"`
	Count int    `json:"qty"`
	Sum   int    `json:"sum"`
}

type Order struct {
	Order_ID        int            `gorm:"primary_key" json:"order_id"`
	User_ID         int            `json:"user_id"`
	Cafe_Id         int            `json:"cafeId"`
	Order_date      string         `json:"date"`
	Cost            int            `json:"cost"`
	Order_list      datatypes.JSON `json:"-"`
	Positions       []Position     `gorm:"-" json:"cart"`
	Address         string         `json:"address"`
	Status_accepted bool
	Status_sent     bool
	Status_canceled bool
}
