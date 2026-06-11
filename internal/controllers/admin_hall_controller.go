package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHallController struct {
	services *services.Services
}

func NewAdminHallController(s *services.Services) *AdminHallController {
	return &AdminHallController{services: s}
}

func (ctl *AdminHallController) ListPackages(c *gin.Context) {
	packages, err := ctl.services.Halls.ListPackagesAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch hall packages.", err.Error())
		return
	}
	utils.OK(c, "Hall packages fetched.", packages)
}

func (ctl *AdminHallController) CreatePackage(c *gin.Context) {
	var input services.HallPackageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid package data.", err.Error())
		return
	}
	pkg, err := ctl.services.Halls.CreatePackage(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to create package.", err.Error())
		return
	}
	utils.Created(c, "Hall package created.", pkg)
}

func (ctl *AdminHallController) UpdatePackage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input services.HallPackageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid package data.", err.Error())
		return
	}
	pkg, err := ctl.services.Halls.UpdatePackage(c.Request.Context(), id, input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update package.", err.Error())
		return
	}
	utils.OK(c, "Hall package updated.", pkg)
}

func (ctl *AdminHallController) DeletePackage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.Halls.DeletePackage(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete package.", err.Error())
		return
	}
	utils.OK(c, "Hall package deleted.", nil)
}

func (ctl *AdminHallController) ListEnquiries(c *gin.Context) {
	page := utils.QueryInt(c, "page", 1)
	limit := utils.QueryInt(c, "limit", 20)
	enquiries, total, err := ctl.services.Halls.ListEnquiries(c.Request.Context(), page, limit)
	if err != nil {
		utils.InternalError(c, "Failed to fetch enquiries.", err.Error())
		return
	}
	utils.OK(c, "Enquiries fetched.", gin.H{"items": enquiries, "total": total})
}

func (ctl *AdminHallController) UpdateEnquiryStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid status.", err.Error())
		return
	}
	enquiry, err := ctl.services.Halls.UpdateEnquiryStatus(c.Request.Context(), id, input.Status, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update enquiry status.", err.Error())
		return
	}
	utils.OK(c, "Enquiry status updated.", enquiry)
}
