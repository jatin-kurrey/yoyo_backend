package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type OfferController struct {
	services *services.Services
}

func NewOfferController(s *services.Services) *OfferController {
	return &OfferController{services: s}
}

func (ctl *OfferController) ListActive(c *gin.Context) {
	offers, err := ctl.services.Offers.ListActive(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to fetch active offers.", err.Error())
		return
	}
	utils.OK(c, "Active offers fetched.", offers)
}
