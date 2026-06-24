package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PMSAccountsController struct {
	svc *services.PMSAccountsService
}

func NewPMSAccountsController(svc *services.PMSAccountsService) *PMSAccountsController {
	return &PMSAccountsController{svc: svc}
}

func (ctl *PMSAccountsController) ListTransactions(c *gin.Context) {
	txns, err := ctl.svc.ListTransactions(c.Request.Context(), c.Query("type"), c.Query("status"))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Transactions loaded.", txns)
}

func (ctl *PMSAccountsController) CreateTransaction(c *gin.Context) {
	var input services.CreateTransactionInput
	if !bindAndValidate(c, &input) {
		return
	}
	tx, err := ctl.svc.CreateTransaction(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Transaction created.", tx)
}

func (ctl *PMSAccountsController) DeleteTransaction(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := ctl.svc.DeleteTransaction(c.Request.Context(), id); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Transaction deleted.", nil)
}

func (ctl *PMSAccountsController) GetSettings(c *gin.Context) {
	settings, err := ctl.svc.GetSettings(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Settings loaded.", settings)
}

func (ctl *PMSAccountsController) UpsertSetting(c *gin.Context) {
	var body struct {
		Key   string `json:"key" validate:"required"`
		Value string `json:"value" validate:"required"`
	}
	if !bindAndValidate(c, &body) {
		return
	}
	if err := ctl.svc.UpsertSetting(c.Request.Context(), body.Key, body.Value); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Setting saved.", nil)
}

func (ctl *PMSAccountsController) ListRateOverrides(c *gin.Context) {
	var catID *uuid.UUID
	if id := c.Query("category_id"); id != "" {
		parsed, err := uuid.Parse(id)
		if err == nil {
			catID = &parsed
		}
	}
	overrides, err := ctl.svc.ListRateOverrides(c.Request.Context(), catID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Rate overrides loaded.", overrides)
}

func (ctl *PMSAccountsController) SetRateOverride(c *gin.Context) {
	var input services.RateOverrideInput
	if !bindAndValidate(c, &input) {
		return
	}
	override, err := ctl.svc.SetRateOverride(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Rate override saved.", override)
}

func (ctl *PMSAccountsController) ClearRateOverride(c *gin.Context) {
	var body struct {
		CategoryID uuid.UUID `json:"category_id" validate:"required"`
		Date       string    `json:"date" validate:"required"`
		Plan       string    `json:"plan" validate:"required"`
	}
	if !bindAndValidate(c, &body) {
		return
	}
	if err := ctl.svc.ClearRateOverride(c.Request.Context(), body.CategoryID, body.Date, body.Plan); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Rate override cleared.", nil)
}
