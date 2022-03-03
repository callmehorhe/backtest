package repository

import (
	"encoding/json"
	"time"

	serv "github.com/callmehorhe/backtest"
	"gorm.io/datatypes"
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

	j, _ := json.Marshal(order.Positions)
	order.Order_list = datatypes.JSON(j)
	order.Order_date = time.Now().Format("2006-01-02 15:04:05")

	row := r.db.Table("orders").Create(&order)
	if err := row.Error; err != nil {
		return 0
	}
	row.Select("order_id").Scan(&id)
	return id
}

func (r *OrderPostgres) UpdateOrder(order serv.Order) serv.Order {
	id := order.Order_ID
	updatedOrder := serv.Order{
		User_ID:    order.User_ID,
		Cafe_Id:    order.Cafe_Id,
		Order_date: order.Order_date,
		Cost:       order.Cost,
		Status:     order.Status,
		Address:    order.Address,
	}
	row := r.db.Table("orders").Where("order_id=?", id).Updates(&updatedOrder)
	row.Take(&updatedOrder)
	return updatedOrder
}

func (r *OrderPostgres) GetOrderByID(id int) serv.Order {
	var order serv.Order
	r.db.Table("orders").Where("order_id=?", id).Take(&order)
	return order
}

func (r *OrderPostgres) GetOrdersByUser(id int) []serv.Order {
	var orders []serv.Order
	r.db.Table("orders").Where("user_id=?", id).Order("order_id desc").Find(&orders)
	return orders
}

func (r *OrderPostgres) GetCafeNameByID(id int) string {
	var name string
	r.db.Table("cafes").Select("name").Where("id_cafe=?", id).Take(&name)
	return name
}
