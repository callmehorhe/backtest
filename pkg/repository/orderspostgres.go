package repository

import (
	"time"

	serv "github.com/callmehorhe/backtest"
	//"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderPostgres struct {
	db *gorm.DB
}

func NewOrderPostgres(db *gorm.DB) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}

func (r *OrderPostgres) CreateOrder(order serv.Order) int {
	var id int
	var list []int
	for i := range order.Positions {
		list = append(list, order.Positions[i].ID)
	}
	order.Order_list = list
	order.Order_date = time.Now().Format("2006-01-02 15:04:05")
	
	row := r.db.Table("orders").Create(&order)
	if err := row.Error; err != nil {
		return 0
	}
	row.Select("order_id").Scan(&id)
	return id
}

func (r *OrderPostgres) UpdateOrder(order serv.Order) {
	id := order.Order_ID
	var list []int
	for i := range order.Positions {
		list = append(list, order.Positions[i].ID)
	}

	r.db.Table("orders").Where("order_id=?", id).Updates(&serv.Order{
		User_ID:         order.User_ID,
		Cafe_Id:         order.Cafe_Id,
		Order_date:      order.Order_date,
		Cost:            order.Cost,
		Order_list:      list,
		Address:         order.Address,
		Status_accepted: order.Status_accepted,
		Status_sent:     order.Status_sent,
		Status_canceled: order.Status_canceled,
	})
}

func (r *OrderPostgres) GetOrderByID(id int) serv.Order {
	var order serv.Order
	r.db.Table("orders").Where("order_id=?", id).Take(&order)
	return order
}
