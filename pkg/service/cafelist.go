package service

import (
	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
)

type CafeService struct {
	repo repository.CafeList
}

func NewCafeService(repo repository.CafeList) *CafeService {
	return &CafeService{
		repo: repo,
	}
}

func (s *CafeService) GetCafeList() []serv.Cafe {
	return s.repo.GetCafeList()
} 

func (s *CafeService) GetMenuByCafeID(id int) []serv.Menu {
	return s.repo.GetMenuByCafeID(id)
}

func (s *CafeService) GetCafeByID(id int) serv.Cafe {
	return s.repo.GetCafeByID(id)
}

func (s *CafeService) GetCategoriesByCafeID(id int) []string {
	return s.repo.GetCategoriesByCafeID(id)
}
