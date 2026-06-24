package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSSystemController struct {
	svc *services.PMSSystemService
}

func NewPMSSystemController(svc *services.PMSSystemService) *PMSSystemController {
	return &PMSSystemController{svc: svc}
}

func (ctl *PMSSystemController) Stats(c *gin.Context) {
	stats, err := ctl.svc.GetStats(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "System stats loaded.", stats)
}

func (ctl *PMSSystemController) Backup(c *gin.Context) {
	data, err := ctl.svc.Export(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=pms-backup.json")
	utils.OK(c, "Backup exported.", data)
}

func (ctl *PMSSystemController) Restore(c *gin.Context) {
	var data services.BackupData
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.BadRequest(c, "Invalid backup file.", nil)
		return
	}
	if data.Version == "" {
		utils.BadRequest(c, "Invalid backup format.", nil)
		return
	}
	if err := ctl.svc.Import(c.Request.Context(), &data); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Backup restored successfully.", nil)
}

func (ctl *PMSSystemController) Reset(c *gin.Context) {
	if err := ctl.svc.ResetAll(c.Request.Context()); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "System has been reset. All PMS data cleared.", nil)
}
