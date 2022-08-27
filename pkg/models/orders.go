package models

import "gorm.io/datatypes"

type Order struct {
	Order_ID   int            `json:"order_id"      gorm:"primary_key;type:serial"`
	User_ID    int            `json:"userId"        gorm:"type:integer"`
	Cafe_Id    int            `json:"cafeId"        gorm:"type:integer"`
	Cafe_Name  string         `json:"cafeName"      gorm:"type:varchar(255)"`
	Order_date string         `json:"date"          gorm:"type:timestamp with time zone"`
	Cost       int            `json:"cost"          gorm:"type:integer"`
	Order_list datatypes.JSON `json:"-"             gorm:"type:jsonb"`
	Positions  []Position     `json:"cart"          gorm:"-"`
	Address    string         `json:"address"       gorm:"type:varchar(255)"`
	Phone      string         `json:"phone"         gorm:"type:varchar(255)"`
	Status     string         `json:"status"        gorm:"type:varchar(255)"`
	Driver     int64          `json:"-"             gorm:"type:bigint"`
}
