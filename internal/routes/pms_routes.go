package routes

import (
	"yoyo-server/internal/config"
	"yoyo-server/internal/controllers"
	"yoyo-server/internal/middleware"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupPMSRoutes(router *gin.Engine, cfg *config.Config, repos *repositories.Repositories, pmsBookingCtl *controllers.PMSBookingController, pmsRoomCtl *controllers.PMSRoomController, pmsFolioCtl *controllers.PMSFolioController, pmsPOSCtl *controllers.PMSPOSController, pmsHKCtl *controllers.PMSHKController, pmsPricingCtl *controllers.PMSPricingController, pmsDashboardCtl *controllers.PMSDashboardController) {
	pms := router.Group("/api/pms")
	pms.Use(middleware.AdminAuth(cfg, repos.AdminUsers))

	// ── Staff-level routes (bookings, folio, POS, housekeeping, room view) ──
	staff := pms.Group("")
	staff.Use(middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin, models.RoleStaff))

	// Dashboard
	staff.GET("/dashboard/stats", pmsDashboardCtl.Stats)

	// Rooms
	staff.GET("/rooms", pmsRoomCtl.List)

	// Bookings
	staff.GET("/bookings", pmsBookingCtl.List)
	staff.POST("/bookings", pmsBookingCtl.Create)
	staff.GET("/bookings/:id", pmsBookingCtl.Get)
	staff.PATCH("/bookings/:id/check-in", pmsBookingCtl.CheckIn)
	staff.PATCH("/bookings/:id/check-out", pmsBookingCtl.CheckOut)
	staff.DELETE("/bookings/:id", pmsBookingCtl.Cancel)

	// Folio
	staff.GET("/bookings/:id/folio", pmsFolioCtl.GetFolio)
	staff.POST("/bookings/:id/folio", pmsFolioCtl.AddEntry)
	staff.POST("/bookings/:id/payments", pmsFolioCtl.AddPayment)

	// POS
	staff.GET("/pos/tables", pmsPOSCtl.ListTables)
	staff.POST("/pos/tables/:id/occupy", pmsPOSCtl.OccupyTable)
	staff.POST("/pos/tables/:id/kot", pmsPOSCtl.AddKOT)
	staff.POST("/pos/tables/:id/bill", pmsPOSCtl.GenerateBill)
	staff.POST("/pos/tables/:id/vacate", pmsPOSCtl.VacateTable)
	staff.POST("/pos/tables/:id/move-to-room", pmsPOSCtl.MoveToRoom)
	staff.GET("/pos/tables/:id/kots", pmsPOSCtl.GetKOTs)

	// Housekeeping
	staff.GET("/housekeeping/tasks", pmsHKCtl.ListTasks)
	staff.POST("/housekeeping/tasks", pmsHKCtl.CreateTask)
	staff.PATCH("/housekeeping/tasks/:id", pmsHKCtl.UpdateTaskStatus)
	staff.PATCH("/housekeeping/rooms/:id/clean", pmsHKCtl.SetRoomClean)
	staff.PATCH("/housekeeping/rooms/:id/dirty", pmsHKCtl.SetRoomDirty)
	staff.PATCH("/housekeeping/rooms/:id/ooo", pmsHKCtl.SetRoomOOO)
	staff.PATCH("/housekeeping/rooms/:id/available", pmsHKCtl.SetRoomAvailable)

	// Pricing (read-only for staff)
	staff.GET("/categories", pmsPricingCtl.ListCategories)

	// ── Admin-level routes (pricing updates, reports, config) ──
	admin := pms.Group("")
	admin.Use(middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin))
	admin.PATCH("/categories/:id/rates", pmsPricingCtl.UpdateRates)
}
