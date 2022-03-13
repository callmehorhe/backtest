package serv

import "gorm.io/datatypes"

type User struct {
	Id_User  int    `json:"-"        gorm:"type:serial;primary_key"`
	Name     string `json:"name"     gorm:"type:varchar(255);not null"`
	Email    string `json:"email"    gorm:"type:varchar(255);unique;not null"`
	Phone    string `json:"phone"    gorm:"type:varchar(255);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}

type Cafe struct {
	Id_Cafe int    `json:"id"      gorm:"primary_key;type:serial"`
	Name    string `json:"name"    gorm:"type:varchar(255)"`
	Phone   string `json:"phone"   gorm:"type:varchar(255)"`
	Image   string `json:"img"     gorm:"type:varchar(255)"`
	Address string `json:"address" gorm:"type:varchar(255)"`
	Chat_ID int64  `json:"-"       gorm:"type:bigint"`
}

type Menu struct {
	Id_Menu  int    `json:"id"          gorm:"primary_key;type:serial"`
	Name     string `json:"productName" gorm:"type:varchar(255)"`
	Image    string `json:"img"         gorm:"type:varchar(255)"`
	Price    int    `json:"cost"        gorm:"type:integer"`
	Category string `json:"category"    gorm:"type:varchar(255)"`
}

type Category struct {
	Category_Name string `json:"name"`
	Menu_List     []Menu `json:"cart"`
}

type CafeAndMenu struct {
	Cafe_Name  string   `json:"name"`
	Categories []string `json:"categories"`
	Menu       []Menu   `json:"menu"`
}

type Position struct {
	ID    int    `json:"productId"`
	Name  string `json:"productName"`
	Price int    `json:"cost"`
	Count int    `json:"qty"`
	Sum   int    `json:"sum"`
}

type Order struct {
	Order_ID   int            `json:"order_id" gorm:"primary_key;type:serial"`
	User_ID    int            `json:"userId"   gorm:"type:integer"`
	Cafe_Id    int            `json:"cafeId"   gorm:"type:integer"`
	Cafe_Name  string         `json:"cafeName" gorm:"-"`
	Order_date string         `json:"date"     gorm:"type:timestamp with time zone"`
	Cost       int            `json:"cost"     gorm:"type:integer"`
	Order_list datatypes.JSON `json:"-"        gorm:"type:jsonb"`
	Positions  []Position     `json:"cart"     gorm:"-"`
	Address    string         `json:"address"  gorm:"type:varchar(255)"`
	Phone      string         `json:"phone"    gorm:"type:varchar(255)"`
	Status     string         `json:"status"   gorm:"type:varchar(255)"`
}
