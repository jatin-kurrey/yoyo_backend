package controllers

import (
	"encoding/json"
	"io"

	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	razorpay *services.RazorpayService
	bookings *services.BookingService
}

func NewWebhookController(razorpay *services.RazorpayService, bookings *services.BookingService) *WebhookController {
	return &WebhookController{razorpay: razorpay, bookings: bookings}
}

func (ctl *WebhookController) Razorpay(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.BadRequest(c, "Invalid webhook body.", nil)
		return
	}

	if !ctl.razorpay.VerifyWebhookSignature(body, c.GetHeader("X-Razorpay-Signature")) {
		utils.BadRequest(c, "Invalid webhook signature.", nil)
		return
	}

	var event struct {
		Event   string `json:"event"`
		Payload struct {
			Payment struct {
				Entity struct {
					ID      string `json:"id"`
					OrderID string `json:"order_id"`
					Status  string `json:"status"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(body, &event); err == nil {
		if event.Event == "payment.failed" || event.Payload.Payment.Entity.Status == "failed" {
			_ = ctl.bookings.MarkPaymentFailedByOrder(c.Request.Context(), event.Payload.Payment.Entity.OrderID, event.Payload.Payment.Entity.ID)
		}
	}

	utils.OK(c, "Webhook accepted.", nil)
}
