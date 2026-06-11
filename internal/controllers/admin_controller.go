package controllers

import (
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminController struct {
	tickets   *services.TicketService
	bookings  *services.BookingService
	contacts  *services.ContactService
	settings  *services.SettingsService
	dashboard *services.DashboardService
	users     *services.AdminUserService
	audit     *services.AuditService
	uploads   *services.UploadService
	hero      *services.HeroSlideService
	content   *services.ContentService
}

func NewAdminController(s *services.Services) *AdminController {
	return &AdminController{
		tickets:   s.Tickets,
		bookings:  s.Bookings,
		contacts:  s.Contacts,
		settings:  s.Settings,
		dashboard: s.Dashboard,
		users:     s.Users,
		audit:     s.Audit,
		uploads:   s.Uploads,
		hero:      s.HeroSlides,
		content:   s.Content,
	}
}

func (ctl *AdminController) DashboardStats(c *gin.Context) {
	stats, err := ctl.dashboard.Stats(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Dashboard stats loaded.", stats)
}

func (ctl *AdminController) ListTickets(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	tickets, total, err := ctl.tickets.ListAdmin(c.Request.Context(), repositories.TicketFilter{
		Search: c.Query("search"),
		Status: c.Query("status"),
	}, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Tickets loaded.", tickets, page, limit, total)
}

func (ctl *AdminController) CreateTicket(c *gin.Context) {
	var input services.TicketInput
	if !bindAndValidate(c, &input) {
		return
	}
	ticket, err := ctl.tickets.Create(c.Request.Context(), input, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Ticket created.", ticket)
}

func (ctl *AdminController) GetTicket(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	ticket, err := ctl.tickets.FindAdminByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Ticket loaded.", ticket)
}

func (ctl *AdminController) UpdateTicket(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var input services.TicketInput
	if !bindAndValidate(c, &input) {
		return
	}
	ticket, err := ctl.tickets.Update(c.Request.Context(), id, input, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Ticket updated.", ticket)
}

func (ctl *AdminController) DeleteTicket(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := ctl.tickets.Delete(c.Request.Context(), id, currentAdminID(c), c.ClientIP()); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Ticket deleted.", nil)
}

func (ctl *AdminController) ToggleTicketStatus(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	ticket, err := ctl.tickets.ToggleStatus(c.Request.Context(), id, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Ticket status updated.", ticket)
}

func (ctl *AdminController) ListBookings(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	bookings, total, err := ctl.bookings.ListAdmin(c.Request.Context(), repositories.BookingFilter{
		Search:        c.Query("search"),
		Status:        c.Query("status"),
		PaymentStatus: c.Query("payment_status"),
		DateFrom:      parseDateQuery(c.Query("date_from")),
		DateTo:        parseDateQuery(c.Query("date_to")),
	}, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Bookings loaded.", bookings, page, limit, total)
}

func (ctl *AdminController) GetBooking(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	booking, err := ctl.bookings.FindAdminByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Booking loaded.", booking)
}

func (ctl *AdminController) UpdateBookingStatus(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var request struct {
		Status string `json:"status" validate:"required"`
	}
	if !bindAndValidate(c, &request) {
		return
	}
	status, ok := validBookingStatus(request.Status)
	if !ok {
		utils.BadRequest(c, "Invalid booking status.", nil)
		return
	}
	booking, err := ctl.bookings.UpdateStatus(c.Request.Context(), id, status, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Booking status updated.", booking)
}

func (ctl *AdminController) ListMessages(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	messages, total, err := ctl.contacts.ListAdmin(c.Request.Context(), repositories.ContactMessageFilter{
		Search: c.Query("search"),
		Status: c.Query("status"),
	}, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Messages loaded.", messages, page, limit, total)
}

func (ctl *AdminController) GetMessage(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	message, err := ctl.contacts.FindAdminByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Message loaded.", message)
}

func (ctl *AdminController) UpdateMessageStatus(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var request struct {
		Status string `json:"status" validate:"required"`
	}
	if !bindAndValidate(c, &request) {
		return
	}
	status, ok := validMessageStatus(request.Status)
	if !ok {
		utils.BadRequest(c, "Invalid message status.", nil)
		return
	}
	message, err := ctl.contacts.UpdateStatus(c.Request.Context(), id, status, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Message status updated.", message)
}

func (ctl *AdminController) DeleteMessage(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := ctl.contacts.Delete(c.Request.Context(), id, currentAdminID(c), c.ClientIP()); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Message deleted.", nil)
}

func (ctl *AdminController) GetSettings(c *gin.Context) {
	setting, err := ctl.settings.Get(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Settings loaded.", setting)
}

func (ctl *AdminController) UpdateSettings(c *gin.Context) {
	var input services.SettingsInput
	if !bindAndValidate(c, &input) {
		return
	}
	setting, err := ctl.settings.Update(c.Request.Context(), input, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Settings updated.", setting)
}

func (ctl *AdminController) ListUsers(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	users, total, err := ctl.users.List(c.Request.Context(), c.Query("search"), page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Admin users loaded.", users, page, limit, total)
}

func (ctl *AdminController) CreateUser(c *gin.Context) {
	var input services.AdminUserInput
	if !bindAndValidate(c, &input) {
		return
	}
	if input.Password == "" {
		utils.BadRequest(c, "Password is required for new admin users.", nil)
		return
	}
	user, err := ctl.users.Create(c.Request.Context(), input, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Admin user created.", user)
}

func (ctl *AdminController) UpdateUser(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var input services.AdminUserInput
	if !bindAndValidate(c, &input) {
		return
	}
	user, err := ctl.users.Update(c.Request.Context(), id, input, currentAdminID(c), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Admin user updated.", user)
}

func (ctl *AdminController) DeleteUser(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := ctl.users.Delete(c.Request.Context(), id, currentAdminID(c), c.ClientIP()); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Admin user deactivated.", nil)
}

func (ctl *AdminController) ListAuditLogs(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	var adminID *uuid.UUID
	if c.Query("admin_id") != "" {
		parsed, err := uuid.Parse(c.Query("admin_id"))
		if err != nil {
			utils.BadRequest(c, "Invalid admin filter.", nil)
			return
		}
		adminID = &parsed
	}
	logs, total, err := ctl.audit.List(c.Request.Context(), repositories.AuditLogFilter{
		Action:  c.Query("action"),
		Module:  c.Query("module"),
		AdminID: adminID,
	}, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Audit logs loaded.", logs, page, limit, total)
}

func (ctl *AdminController) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "A file is required.", nil)
		return
	}
	folder := c.DefaultQuery("folder", "general")
	asset, err := ctl.uploads.Save(c.Request.Context(), *currentAdminID(c), fileHeader, folder, c.ClientIP())
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}
	utils.Created(c, "File uploaded.", asset)
}

func (ctl *AdminController) ListMedia(c *gin.Context) {
	page, limit, _ := utils.ParsePagination(c)
	assets, total, err := ctl.uploads.List(c.Request.Context(), page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Paginated(c, "Media library loaded.", assets, page, limit, total)
}

func (ctl *AdminController) DeleteMedia(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := ctl.uploads.Delete(c.Request.Context(), *currentAdminID(c), id, c.ClientIP()); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Media asset deleted.", nil)
}

func roleCanManageUsers(role string) bool {
	return role == string(models.RoleSuperAdmin) || role == string(models.RoleAdmin)
}

func (ctl *AdminController) ListHeroSlides(c *gin.Context) {
	slides, err := ctl.hero.ListAdmin(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Hero slides loaded.", slides)
}

func (ctl *AdminController) CreateHeroSlide(c *gin.Context) {
	var input services.HeroSlideInput
	if !bindAndValidate(c, &input) {
		return
	}
	slide, err := ctl.hero.Create(c.Request.Context(), input, currentAdminID(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.Created(c, "Hero slide created.", slide)
}

func (ctl *AdminController) UpdateHeroSlide(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	var input services.HeroSlideInput
	if !bindAndValidate(c, &input) {
		return
	}
	slide, err := ctl.hero.Update(c.Request.Context(), id, input, currentAdminID(c))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Hero slide updated.", slide)
}

func (ctl *AdminController) DeleteHeroSlide(c *gin.Context) {
	id, ok := uuidParam(c, "id")
	if !ok {
		return
	}
	if err := ctl.hero.Delete(c.Request.Context(), id, currentAdminID(c)); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Hero slide deleted.", nil)
}

func (ctl *AdminController) ListContent(c *gin.Context) {
	pages, err := ctl.content.List(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Content pages loaded.", pages)
}

func (ctl *AdminController) GetContent(c *gin.Context) {
	page, err := ctl.content.AdminFindBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Content page loaded.", page)
}

func (ctl *AdminController) UpdateContent(c *gin.Context) {
	var input services.UpdateContentInput
	if !bindAndValidate(c, &input) {
		return
	}
	page, err := ctl.content.Update(c.Request.Context(), *currentAdminID(c), c.Param("slug"), input, c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.OK(c, "Content page updated.", page)
}
