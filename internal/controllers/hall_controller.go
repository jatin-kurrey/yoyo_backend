package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type HallController struct {
	services *services.Services
}

func NewHallController(s *services.Services) *HallController {
	return &HallController{services: s}
}

func (ctl *HallController) ListPackages(c *gin.Context) {
	packages, err := ctl.services.Halls.ListPackagesPublic(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch hall packages.", err.Error())
		return
	}
	utils.OK(c, "Hall packages fetched.", packages)
}

func (ctl *HallController) SubmitEnquiry(c *gin.Context) {
	var input services.HallEnquiryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid enquiry data.", err.Error())
		return
	}
	enquiry, err := ctl.services.Halls.CreateEnquiry(c.Request.Context(), input)
	if err != nil {
		utils.InternalError(c, "Failed to submit enquiry.", err.Error())
		return
	}
	utils.Created(c, "Enquiry submitted successfully.", enquiry)
}
