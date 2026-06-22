package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PMSPOSController struct {
	svc *services.PMSPOSService
}

func NewPMSPOSController(svc *services.PMSPOSService) *PMSPOSController {
	return &PMSPOSController{svc: svc}
}

func (ctl *PMSPOSController) ListTables(c *gin.Context) {
	tables, err := ctl.svc.ListTables(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Tables loaded.", tables)
}

func (ctl *PMSPOSController) OccupyTable(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		GuestName string `json:"guest_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "Invalid request body.", nil)
		return
	}
	table, err := ctl.svc.OccupyTable(c.Request.Context(), id, body.GuestName)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Table occupied.", table)
}

func (ctl *PMSPOSController) AddKOT(c *gin.Context) {
	var input services.AddKOTInput
	if !bindAndValidate(c, &input) {
		return
	}
	order, err := ctl.svc.AddKOT(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "KOT sent.", order)
}

func (ctl *PMSPOSController) GenerateBill(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	table, err := ctl.svc.GenerateBill(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Bill generated.", table)
}

func (ctl *PMSPOSController) VacateTable(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	table, err := ctl.svc.VacateTable(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Table vacated.", table)
}

func (ctl *PMSPOSController) MoveToRoom(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		BookingID uuid.UUID `json:"booking_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "Invalid request body.", nil)
		return
	}
	if err := ctl.svc.MoveToRoomFolio(c.Request.Context(), id, body.BookingID); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Bill moved to room folio.", nil)
}

func (ctl *PMSPOSController) GetKOTs(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	orders, err := ctl.svc.GetKOTs(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "KOTs loaded.", orders)
}
