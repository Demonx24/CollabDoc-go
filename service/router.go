package service

import (
	"CollabDoc-go/model/database"
	"gorm.io/gorm"
	"net/http"
)

type RouterService struct{}

func (service *RouterService) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func fetchAllRoutes(db *gorm.DB) ([]database.RouteMenu, error) {
	var menus []database.RouteMenu
	if err := db.Order("parent_id, sort_order").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
