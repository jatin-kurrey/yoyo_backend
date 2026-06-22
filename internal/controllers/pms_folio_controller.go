package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSFolioController struct {
	svc *services.PMSBookingsService
}

func NewPMSFolioController(svc *services.PMSBookingsService) *PMSFolioController {
	return &PMSFolioController{svc: svc}
}

func (ctl *PMSFolioController) GetFolio(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	entries, payments, err := ctl.svc.GetFolio(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Folio loaded.", gin.H{"entries": entries, "payments": payments})
}

func (ctl *PMSFolioController) AddEntry(c *gin.Context) {
	var input services.AddFolioEntryInput
	if !bindAndValidate(c, &input) {
		return
	}
	entry, err := ctl.svc.AddFolioEntry(c.Request.Context(), input, currentAdminID(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Entry added to folio.", entry)
}

func (ctl *PMSFolioController) AddPayment(c *gin.Context) {
	var input services.AddPaymentInput
	if !bindAndValidate(c, &input) {
		return
	}
	payment, err := ctl.svc.AddPayment(c.Request.Context(), input, currentAdminID(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Payment recorded.", payment)
}
