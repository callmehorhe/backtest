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