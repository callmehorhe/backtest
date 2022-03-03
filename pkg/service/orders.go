package service

import (
	"encoding/json"
	"log"

	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
)

type OrderService struct {
	repo repository.Orders
}

func NewOrderService(repo repository.Orders) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) GetOrdersByUser(id int) []serv.Order {
	orders := s.repo.GetOrdersByUser(id)
	for i := 0; i < len(orders); i++ {
		orders[i].Cafe_Name = s.repo.GetCafeNameByID(orders[i].Cafe_Id)
		pos := []serv.Position{}
		if err := json.Unmarshal(orders[i].Order_list, &pos); err != nil {
			log.Fatal(err)
		}
		orders[i].Positions = pos
	}
	return orders
}
