package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSRoomController struct {
	pricingSvc *services.PMSPricingService
}

func NewPMSRoomController(pricingSvc *services.PMSPricingService) *PMSRoomController {
	return &PMSRoomController{pricingSvc: pricingSvc}
}

func (ctl *PMSRoomController) List(c *gin.Context) {
	rooms, err := ctl.pricingSvc.ListRooms(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Rooms loaded.", rooms)
}
