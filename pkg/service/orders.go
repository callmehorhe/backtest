package service

import (
	"encoding/json"
	"math"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/sirupsen/logrus"
)

type OrderService struct {
	repo repository.Repository
}

func NewOrderService(repo repository.Repository) *OrderService {
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
			logrus.Error(err)
		}
		orders[i].Positions = pos
		if orders[i].Driver_Id != 0 {
			var err error
			orders[i].Driver, err = s.repo.Drivers.GetDriverById(orders[i].Driver_Id)
			if err != nil {
				logrus.Error("cant add driver: %v", err)
			}
		}

	}
	return orders
}

func (s *OrderService) GetPagesCount(id int) int {
	return int(math.Ceil(float64(s.repo.GetOrdersCount(id)) / 10))
}
