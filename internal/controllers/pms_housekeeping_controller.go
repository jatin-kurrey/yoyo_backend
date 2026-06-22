package controllers

import (
	"yoyo-server/internal/models"
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PMSHKController struct {
	svc *services.PMSHKService
}

func NewPMSHKController(svc *services.PMSHKService) *PMSHKController {
	return &PMSHKController{svc: svc}
}

func (ctl *PMSHKController) ListTasks(c *gin.Context) {
	tasks, err := ctl.svc.ListTasks(c.Request.Context(), c.Query("status"))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Tasks loaded.", tasks)
}

func (ctl *PMSHKController) CreateTask(c *gin.Context) {
	var input services.CreateHKTaskInput
	if !bindAndValidate(c, &input) {
		return
	}
	task, err := ctl.svc.CreateTask(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Task created.", task)
}

func (ctl *PMSHKController) UpdateTaskStatus(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		Status models.HKTaskStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "Invalid request body.", nil)
		return
	}
	task, err := ctl.svc.UpdateTaskStatus(c.Request.Context(), id, body.Status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Task updated.", task)
}

func (ctl *PMSHKController) SetRoomClean(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	room, err := ctl.svc.SetRoomCleanStatus(c.Request.Context(), id, models.CleanClean)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Room marked clean.", room)
}

func (ctl *PMSHKController) SetRoomDirty(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	room, err := ctl.svc.SetRoomCleanStatus(c.Request.Context(), id, models.CleanDirty)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Room marked dirty.", room)
}

func (ctl *PMSHKController) SetRoomOOO(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		body.Reason = "Maintenance"
	}
	room, err := ctl.svc.SetRoomOOO(c.Request.Context(), id, body.Reason)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Room marked out of order.", room)
}

func (ctl *PMSHKController) SetRoomAvailable(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	room, err := ctl.svc.SetRoomAvailable(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Room marked available.", room)
}
