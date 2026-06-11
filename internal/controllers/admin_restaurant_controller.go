package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminRestaurantController struct {
	services *services.Services
}

func NewAdminRestaurantController(s *services.Services) *AdminRestaurantController {
	return &AdminRestaurantController{services: s}
}

func (ctl *AdminRestaurantController) ListItems(c *gin.Context) {
	items, err := ctl.services.Restaurant.ListAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch menu items.", err.Error())
		return
	}
	utils.OK(c, "Menu items fetched.", items)
}

func (ctl *AdminRestaurantController) CreateItem(c *gin.Context) {
	var input services.RestaurantItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid menu data.", err.Error())
		return
	}
	item, err := ctl.services.Restaurant.CreateItem(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to create menu item.", err.Error())
		return
	}
	utils.Created(c, "Menu item created.", item)
}

func (ctl *AdminRestaurantController) UpdateItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input services.RestaurantItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid menu data.", err.Error())
		return
	}
	item, err := ctl.services.Restaurant.UpdateItem(c.Request.Context(), id, input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update menu item.", err.Error())
		return
	}
	utils.OK(c, "Menu item updated.", item)
}

func (ctl *AdminRestaurantController) DeleteItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.Restaurant.DeleteItem(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete menu item.", err.Error())
		return
	}
	utils.OK(c, "Menu item deleted.", nil)
}
