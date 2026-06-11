package controllers

import (
	"time"

	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PublicController struct {
	tickets  *services.TicketService
	bookings *services.BookingService
	contacts *services.ContactService
	settings *services.SettingsService
	hero     *services.HeroSlideService
	content  *services.ContentService
}

func NewPublicController(tickets *services.TicketService, bookings *services.BookingService, contacts *services.ContactService, settings *services.SettingsService, hero *services.HeroSlideService, content *services.ContentService) *PublicController {
	return &PublicController{tickets: tickets, bookings: bookings, contacts: contacts, settings: settings, hero: hero, content: content}
}

func (ctl *PublicController) Tickets(c *gin.Context) {
	tickets, err := ctl.tickets.ListPublic(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Tickets loaded.", tickets)
}

func (ctl *PublicController) TicketBySlug(c *gin.Context) {
	ticket, err := ctl.tickets.FindPublicBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Ticket loaded.", ticket)
}

type createOrderRequest struct {
	CustomerName  string `json:"customer_name" validate:"required,min=2,max=140"`
	CustomerEmail string `json:"customer_email" validate:"required,email,max=180"`
	CustomerPhone string `json:"customer_phone" validate:"required,min=8,max=30"`
	TicketID      string `json:"ticket_id" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required,gte=1"`
	VisitDate     string `json:"visit_date" validate:"required"`
}

func (ctl *PublicController) CreateOrder(c *gin.Context) {
	var request createOrderRequest
	if !bindAndValidate(c, &request) {
		return
	}
	ticketID, err := uuid.Parse(request.TicketID)
	if err != nil {
		utils.BadRequest(c, "Invalid ticket selected.", nil)
		return
	}
	visitDate, err := time.Parse("2006-01-02", request.VisitDate)
	if err != nil {
		utils.BadRequest(c, "Visit date must use YYYY-MM-DD format.", nil)
		return
	}
	if visitDate.Before(time.Now().Truncate(24 * time.Hour)) {
		utils.BadRequest(c, "Visit date cannot be in the past.", nil)
		return
	}

	result, err := ctl.bookings.CreateOrder(c.Request.Context(), services.CreateOrderInput{
		CustomerName:  request.CustomerName,
		CustomerEmail: request.CustomerEmail,
		CustomerPhone: request.CustomerPhone,
		TicketID:      ticketID,
		Quantity:      request.Quantity,
		VisitDate:     visitDate,
	})
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Razorpay order created.", result)
}

func (ctl *PublicController) VerifyPayment(c *gin.Context) {
	var input services.VerifyPaymentInput
	if !bindAndValidate(c, &input) {
		return
	}
	booking, err := ctl.bookings.VerifyPayment(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Payment verified and booking confirmed.", booking)
}

func (ctl *PublicController) Contact(c *gin.Context) {
	var input services.ContactInput
	if !bindAndValidate(c, &input) {
		return
	}
	message, err := ctl.contacts.Create(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Message received.", message)
}

func (ctl *PublicController) PublicSettings(c *gin.Context) {
	setting, err := ctl.settings.Get(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Public settings loaded.", services.PublicSettings(setting))
}

func (ctl *PublicController) HeroSlides(c *gin.Context) {
	slides, err := ctl.hero.ListPublic(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Hero slides loaded.", slides)
}

func (ctl *PublicController) GetContent(c *gin.Context) {
	page, err := ctl.content.FindBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Content page loaded.", page)
}
