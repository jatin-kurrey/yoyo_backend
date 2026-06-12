package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminAttractionController struct {
	services *services.Services
}

func NewAdminAttractionController(s *services.Services) *AdminAttractionController {
	return &AdminAttractionController{services: s}
}

func (ctl *AdminAttractionController) List(c *gin.Context) {
	items, err := ctl.services.Attractions.ListAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch attractions.", err.Error())
		return
	}
	utils.OK(c, "Attractions fetched.", items)
}

func (ctl *AdminAttractionController) Create(c *gin.Context) {
	var input services.AttractionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid attraction data.", err.Error())
		return
	}
	item, err := ctl.services.Attractions.Create(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to create attraction.", err.Error())
		return
	}
	utils.Created(c, "Attraction created.", item)
}

func (ctl *AdminAttractionController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input services.AttractionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid attraction data.", err.Error())
		return
	}
	item, err := ctl.services.Attractions.Update(c.Request.Context(), id, input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update attraction.", err.Error())
		return
	}
	utils.OK(c, "Attraction updated.", item)
}

func (ctl *AdminAttractionController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.Attractions.Delete(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete attraction.", err.Error())
		return
	}
	utils.OK(c, "Attraction deleted.", nil)
}
