package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type SEOController struct {
	services *services.Services
}

func NewSEOController(s *services.Services) *SEOController {
	return &SEOController{services: s}
}

func (ctl *SEOController) GetSEO(c *gin.Context) {
	slug := c.Param("slug")
	seo, err := ctl.services.SEO.GetPublic(c.Request.Context(), slug)
	if err != nil {
		utils.NotFound(c, "SEO not found for this page.", err.Error())
		return
	}
	utils.OK(c, "SEO metadata fetched.", seo)
}
