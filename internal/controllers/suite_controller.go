package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type SuiteController struct {
	services *services.Services
}

func NewSuiteController(s *services.Services) *SuiteController {
	return &SuiteController{services: s}
}

func (ctl *SuiteController) List(c *gin.Context) {
	suites, err := ctl.services.Suites.ListPublic(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch suites.", err.Error())
		return
	}
	utils.OK(c, "Suites fetched.", suites)
}

func (ctl *SuiteController) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	suite, err := ctl.services.Suites.FindBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.NotFound(c, "Suite not found.", err.Error())
		return
	}
	utils.OK(c, "Suite details fetched.", suite)
}
