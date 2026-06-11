package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type GalleryController struct {
	services *services.Services
}

func NewGalleryController(s *services.Services) *GalleryController {
	return &GalleryController{services: s}
}

func (ctl *GalleryController) List(c *gin.Context) {
	items, err := ctl.services.Gallery.ListPublic(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch gallery items.", err.Error())
		return
	}
	utils.OK(c, "Gallery items fetched.", items)
}
