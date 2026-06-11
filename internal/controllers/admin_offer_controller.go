package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminOfferController struct {
	services *services.Services
}

func NewAdminOfferController(s *services.Services) *AdminOfferController {
	return &AdminOfferController{services: s}
}

func (ctl *AdminOfferController) List(c *gin.Context) {
	offers, err := ctl.services.Offers.ListAdmin(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch offers.", err.Error())
		return
	}
	utils.OK(c, "Offers fetched.", offers)
}

func (ctl *AdminOfferController) Create(c *gin.Context) {
	var input services.OfferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid offer data.", err.Error())
		return
	}
	offer, err := ctl.services.Offers.Create(c.Request.Context(), input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to create offer.", err.Error())
		return
	}
	utils.Created(c, "Offer created.", offer)
}

func (ctl *AdminOfferController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	var input services.OfferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Invalid offer data.", err.Error())
		return
	}
	offer, err := ctl.services.Offers.Update(c.Request.Context(), id, input, *currentAdminID(c), c.ClientIP())
	if err != nil {
		utils.InternalError(c, "Failed to update offer.", err.Error())
		return
	}
	utils.OK(c, "Offer updated.", offer)
}

func (ctl *AdminOfferController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid ID.", nil)
		return
	}
	if err := ctl.services.Offers.Delete(c.Request.Context(), id, *currentAdminID(c), c.ClientIP()); err != nil {
		utils.InternalError(c, "Failed to delete offer.", err.Error())
		return
	}
	utils.OK(c, "Offer deleted.", nil)
}
