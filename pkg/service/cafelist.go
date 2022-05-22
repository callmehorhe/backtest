package service

import (
	"time"

	serv "github.com/callmehorhe/backtest"
	"github.com/callmehorhe/backtest/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

type CafeService struct {
	repo repository.CafeList
}

func NewCafeService(repo repository.CafeList) *CafeService {
	return &CafeService{
		repo: repo,
	}
}

func (s *CafeService) GetCafe(id int, password string) (serv.Cafe, error) {
	return s.repo.GetCafe(id, password)
}

func (s *CafeService) CafeGenerateToken(id int, password string) (string, error) {
	cafe, err := s.repo.GetCafe(id, password)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		cafe.Id_Cafe,
	})

	return token.SignedString([]byte(signingKey))
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
