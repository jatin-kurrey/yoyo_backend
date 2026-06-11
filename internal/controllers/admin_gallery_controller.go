package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminGalleryController struct {
	services *services.Services
}

func NewAdminGalleryController(s *services.Services) *AdminGalleryController {
	return &AdminGalleryController{services: s}
}

func (ctl *AdminGalleryController) List(c *gin.Context) {
	items, err := ctl.services.Gallery.ListAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch gallery items.", err.Error())
		return
	}
	utils.OK(c, "Gallery items fetched.", items)
}

func (ctl *AdminGalleryController) Create(c *gin.Context) {
	var input services.GalleryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid gallery data.", err.Error())
		return
	}
	item, err := ctl.services.Gallery.Create(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to create gallery item.", err.Error())
		return
	}
	utils.Created(c, "Gallery item created.", item)
}

func (ctl *AdminGalleryController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input services.GalleryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid gallery data.", err.Error())
		return
	}
	item, err := ctl.services.Gallery.Update(c.Request.Context(), id, input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update gallery item.", err.Error())
		return
	}
	utils.OK(c, "Gallery item updated.", item)
}

func (ctl *AdminGalleryController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.Gallery.Delete(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete gallery item.", err.Error())
		return
	}
	utils.OK(c, "Gallery item deleted.", nil)
}
