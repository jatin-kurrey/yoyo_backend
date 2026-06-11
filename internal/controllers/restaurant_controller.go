package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type RestaurantController struct {
	services *services.Services
}

func NewRestaurantController(s *services.Services) *RestaurantController {
	return &RestaurantController{services: s}
}

func (ctl *RestaurantController) ListItems(c *gin.Context) {
	items, err := ctl.services.Restaurant.ListPublic(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch menu items.", err.Error())
		return
	}
	utils.OK(c, "Menu items fetched.", items)
}
