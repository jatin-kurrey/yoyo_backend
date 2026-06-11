package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminSEOController struct {
	services *services.Services
}

func NewAdminSEOController(s *services.Services) *AdminSEOController {
	return &AdminSEOController{services: s}
}

func (ctl *AdminSEOController) List(c *gin.Context) {
	pages, err := ctl.services.SEO.ListAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch SEO pages.", err.Error())
		return
	}
	utils.OK(c, "SEO pages fetched.", pages)
}

func (ctl *AdminSEOController) Save(c *gin.Context) {
	var input services.SEOInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid SEO data.", err.Error())
		return
	}
	page, err := ctl.services.SEO.Save(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to save SEO data.", err.Error())
		return
	}
	utils.OK(c, "SEO data saved.", page)
}

func (ctl *AdminSEOController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	page, err := ctl.services.SEO.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, "SEO page not found.", err.Error())
		return
	}
	utils.OK(c, "SEO page fetched.", page)
}

func (ctl *AdminSEOController) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	page, err := ctl.services.SEO.GetPublic(c.Request.Context(), slug)
	if err != nil {
		utils.NotFound(c, "SEO page not found.", err.Error())
		return
	}
	utils.OK(c, "SEO page fetched.", page)
}

func (ctl *AdminSEOController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.SEO.Delete(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete SEO page.", err.Error())
		return
	}
	utils.OK(c, "SEO page deleted.", nil)
}
