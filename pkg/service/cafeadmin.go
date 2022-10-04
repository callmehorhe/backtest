package service

import (
	"fmt"
	"strings"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/spf13/viper"
)

func (s *CafeService) UpdateCafe(cafe models.Cafe) error {
	if cafe.BaseImage != "" {
		path := strings.ReplaceAll(fmt.Sprintf("%s/%sName.jpg", cafe.Name, cafe.Name), " ", "")
		SaveImage(cafe.BaseImage, path)
		cafe.Image = "http://92.63.104.228" + viper.GetString("port") + "/images/" + path
	}
	return s.repo.UpdateCafe(cafe)
}

func (s *CafeService) UpdateMenu(menu []models.Menu, cafe string) {
	oldMenu := s.repo.GetMenuByCafeID(menu[0].Id_Cafe)
	var del []int
	for _, pos := range menu {
		if pos.BaseImage != "" {
			path := strings.ReplaceAll(fmt.Sprintf("%s/%s.jpg", cafe, pos.Name), " ", "")
			SaveImage(pos.BaseImage, path)
			pos.Image = "http://92.63.104.228" + viper.GetString("port") + "/images/" + path
		}
		if pos.Id_Menu < 0 {
			s.repo.CreatePos(pos)
		} else {
			s.repo.UpdatePos(pos)
		}
	}

	for _, oldId := range oldMenu {
		exist := false
		for _, newID := range menu {
			if newID.Id_Menu == oldId.Id_Menu {
				exist = true
				break
			}
		}
		if !exist {
			del = append(del, oldId.Id_Menu)
		}

	}
	if len(del) > 0 {
		s.repo.DeletePos(del)
	}

}
