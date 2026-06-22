package controllers

import (
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSBookingController struct {
	svc *services.PMSBookingService
}

func NewPMSBookingController(svc *services.PMSBookingService) *PMSBookingController {
	return &PMSBookingController{svc: svc}
}

func (ctl *PMSBookingController) List(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	bookings, total, err := ctl.svc.List(c.Request.Context(), c.Query("search"), c.Query("status"), page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Bookings loaded.", bookings, page, limit, total)
}

func (ctl *PMSBookingController) Get(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	booking, err := ctl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Booking loaded.", booking)
}

func (ctl *PMSBookingController) Create(c *gin.Context) {
	var input services.CreatePMSBookingInput
	if !bindAndValidate(c, &input) {
		return
	}
	booking, err := ctl.svc.Create(c.Request.Context(), input, currentAdminID(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Booking created.", booking)
}

func (ctl *PMSBookingController) CheckIn(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	booking, err := ctl.svc.CheckIn(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Guest checked in.", booking)
}

func (ctl *PMSBookingController) CheckOut(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	booking, err := ctl.svc.CheckOut(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Guest checked out.", booking)
}

func (ctl *PMSBookingController) Cancel(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	booking, err := ctl.svc.Cancel(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Booking cancelled.", booking)
}
