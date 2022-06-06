package service

import (
	"encoding/json"
	"log"
	"math"

	"github.com/callmehorhe/backtest/pkg/models"
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

func (s *OrderService) GetOrdersByUser(id, count int) []models.Order {
	orders := s.repo.GetOrdersByUser(id, count)
	for i := 0; i < len(orders); i++ {
		orders[i].Cafe_Name = s.repo.GetCafeNameByID(orders[i].Cafe_Id)
		pos := []models.Position{}
		if err := json.Unmarshal(orders[i].Order_list, &pos); err != nil {
			log.Fatal(err)
		}
		orders[i].Positions = pos
	}
	return orders
}

func (s *OrderService) GetPagesCount(id int) int {
	return int(math.Ceil(float64(s.repo.GetOrdersCount(id)) / 10))
}
