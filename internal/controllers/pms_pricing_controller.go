package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSPricingController struct {
	svc *services.PMSPricingService
}

func NewPMSPricingController(svc *services.PMSPricingService) *PMSPricingController {
	return &PMSPricingController{svc: svc}
}

func (ctl *PMSPricingController) ListCategories(c *gin.Context) {
	cats, err := ctl.svc.ListCategoriesWithRooms(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Categories loaded.", cats)
}

func (ctl *PMSPricingController) UpdateRates(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		BasePrice int64 `json:"base_price"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "Invalid request body.", nil)
		return
	}
	cat, err := ctl.svc.UpdateRates(c.Request.Context(), id, body.BasePrice)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Rates updated.", cat)
}
