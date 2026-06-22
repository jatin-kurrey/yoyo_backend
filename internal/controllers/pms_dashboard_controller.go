package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSDashboardController struct {
	svc *services.PMSDashboardService
}

func NewPMSDashboardController(svc *services.PMSDashboardService) *PMSDashboardController {
	return &PMSDashboardController{svc: svc}
}

func (ctl *PMSDashboardController) Stats(c *gin.Context) {
	stats, err := ctl.svc.GetStats(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Dashboard stats loaded.", stats)
}
