package controllers

import (
	"errors"
	"net/http"
	"time"

	"yoyo-server/internal/middleware"
	"yoyo-server/internal/models"
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"
	"yoyo-server/internal/validators"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func bindAndValidate(c *gin.Context, payload interface{}) bool {
	if err := c.ShouldBindJSON(payload); err != nil {
		utils.BadRequest(c, "Invalid JSON payload.", err.Error())
		return false
	}
	if errors := validators.Struct(payload); errors != nil {
		utils.BadRequest(c, "Please correct the highlighted fields.", errors)
		return false
	}
	return true
}

func uuidParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		utils.BadRequest(c, "Invalid identifier.", nil)
		return uuid.Nil, false
	}
	return id, true
}

func currentAdminID(c *gin.Context) *uuid.UUID {
	user, ok := middleware.CurrentAdmin(c)
	if !ok {
		return nil
	}
	return &user.ID
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidCredentials):
		utils.Unauthorized(c, "Invalid email or password.")
	case errors.Is(err, services.ErrInactiveAccount):
		utils.Forbidden(c, "This admin account is inactive.")
	case errors.Is(err, services.ErrNotFound):
		utils.NotFound(c, "The requested resource was not found.", nil)
	case errors.Is(err, services.ErrInsufficientStock):
		utils.BadRequest(c, "Selected ticket does not have enough stock.", nil)
	case errors.Is(err, services.ErrInvalidSignature):
		utils.BadRequest(c, "Payment signature verification failed.", nil)
	case errors.Is(err, services.ErrPaymentGatewayDisabled):
		utils.Error(c, http.StatusServiceUnavailable, "Payment gateway is not configured.", nil)
	case errors.Is(err, services.ErrOnlySuperAdmin):
		utils.BadRequest(c, "At least one active super admin must remain.", nil)
	default:
		utils.ServerError(c)
	}
}

func parseDateQuery(value string) *time.Time {
	if value == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil
	}
	return &parsed
}

func validBookingStatus(status string) (models.BookingStatus, bool) {
	switch models.BookingStatus(status) {
	case models.BookingPending, models.BookingConfirmed, models.BookingCancelled, models.BookingRefunded:
		return models.BookingStatus(status), true
	default:
		return "", false
	}
}

func validMessageStatus(status string) (models.MessageStatus, bool) {
	switch models.MessageStatus(status) {
	case models.MessageNew, models.MessageRead, models.MessageReplied, models.MessageArchived:
		return models.MessageStatus(status), true
	default:
		return "", false
	}
}
