package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type AttractionController struct {
	services *services.Services
}

func NewAttractionController(s *services.Services) *AttractionController {
	return &AttractionController{services: s}
}

func (ctl *AttractionController) List(c *gin.Context) {
	items, err := ctl.services.Attractions.ListPublic(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch attractions.", err.Error())
		return
	}
	utils.OK(c, "Attractions fetched.", items)
}
