package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminSuiteController struct {
	services *services.Services
}

func NewAdminSuiteController(s *services.Services) *AdminSuiteController {
	return &AdminSuiteController{services: s}
}

func (ctl *AdminSuiteController) List(c *gin.Context) {
	suites, err := ctl.services.Suites.ListAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch suites.", err.Error())
		return
	}
	utils.OK(c, "Suites fetched.", suites)
}

func (ctl *AdminSuiteController) Create(c *gin.Context) {
	var input services.SuiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid suite data.", err.Error())
		return
	}
	suite, err := ctl.services.Suites.Create(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to create suite.", err.Error())
		return
	}
	utils.Created(c, "Suite created.", suite)
}

func (ctl *AdminSuiteController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input services.SuiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid suite data.", err.Error())
		return
	}
	suite, err := ctl.services.Suites.Update(c.Request.Context(), id, input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update suite.", err.Error())
		return
	}
	utils.OK(c, "Suite updated.", suite)
}

func (ctl *AdminSuiteController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.Suites.Delete(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete suite.", err.Error())
		return
	}
	utils.OK(c, "Suite deleted.", nil)
}
