package models

type Cafe struct {
	Id_Cafe   int    `json:"id"      gorm:"primary_key;type:serial"`
	Name      string `json:"name"    gorm:"type:varchar(255)"`
	Phone     string `json:"phone"   gorm:"type:varchar(255)"`
	Image     string `json:"img"     gorm:"type:varchar(255)"`
	BaseImage string `json:"baseimg" gorm:"-"`
	Address   string `json:"address" gorm:"type:varchar(255)"`
	Chat_ID   int64  `json:"-"       gorm:"type:bigint"`
	Password  string `json:"-"       gorm:"type:varchar(255)"`
}

type Menu struct {
	Id_Menu     int    `json:"id"          gorm:"primary_key;type:serial"`
	Id_Cafe     int    `json:"id_cafe"     gorm:"type:varchar(255)"`
	Name        string `json:"productName" gorm:"type:varchar(255)"`
	Image       string `json:"img"         gorm:"type:varchar(255)"`
	BaseImage   string `json:"baseimg"     gorm:"-"`
	Price       int    `json:"cost"        gorm:"type:integer"`
	Category    string `json:"category"    gorm:"type:varchar(255)"`
	Description string `json:"desc"        gorm:"type:varchar(255)"`
	Weight      int    `json:"weight"      gorm:"type:integer"`
	Avaible     bool   `json:"avaible"     gorm:"type:boolean"`
}

type Category struct {
	Category_Name string `json:"name"`
	Menu_List     []Menu `json:"cart"`
}

type CafeAndMenu struct {
	Cafe       Cafe     `json:"cafe"`
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
