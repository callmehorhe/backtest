package service

import (
	"fmt"

	"github.com/callmehorhe/backtest/pkg/models"
)

func (s *CafeService) UpdateCafe(cafe models.Cafe) error {
	if cafe.BaseImage != "" {
		path := fmt.Sprintf("%s/%s.jpg", cafe.Name, cafe.Name)
		SaveImage(cafe.BaseImage, cafe.Name)
		cafe.Image = path
	}
	return s.repo.UpdateCafe(cafe)
}

func (s *CafeService) UpdateMenu(menu []models.Menu) {
	for _, pos := range menu {
		if pos.Id_Menu < 0 {
			s.repo.CreatePos(pos)
		} else {
			s.repo.UpdatePos(pos)
		}
	}
}
