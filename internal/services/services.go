package services

import (
	"yoyo-server/internal/config"
	"yoyo-server/internal/repositories"

	"gorm.io/gorm"
)

type Services struct {
	Auth      *AuthService
	Tickets   *TicketService
	Bookings  *BookingService
	Contacts  *ContactService
	Settings  *SettingsService
	Dashboard *DashboardService
	Users     *AdminUserService
	Audit     *AuditService
	Razorpay  *RazorpayService
	Uploads   *UploadService
	HeroSlides *HeroSlideService
	Content    *ContentService
	Gallery    *GalleryService
	Restaurant *RestaurantService
	Suites     *SuiteService
	Halls      *HallService
	SEO        *SEOService
	Offers     *OfferService
}

func New(cfg *config.Config, db *gorm.DB, repos *repositories.Repositories) *Services {
	audit := NewAuditService(repos.AuditLogs)
	razorpay := NewRazorpayService(cfg)

	var storage StorageProvider
	if cfg.UploadsStorage == "r2" {
		r2, err := NewR2StorageProvider(cfg)
		if err != nil {
			panic(err) // Critical failure if R2 is requested but misconfigured
		}
		storage = r2
	} else {
		storage = NewLocalStorageProvider(cfg.UploadDir)
	}

	return &Services{
		Auth:      NewAuthService(cfg, repos.AdminUsers),
		Tickets:   NewTicketService(repos.Tickets, audit),
		Bookings:  NewBookingService(db, repos.Tickets, repos.Bookings, razorpay, audit),
		Contacts:  NewContactService(repos.Messages, audit),
		Settings:  NewSettingsService(repos.Settings, audit),
		Dashboard: NewDashboardService(repos.Tickets, repos.Bookings, repos.Messages),
		Users:     NewAdminUserService(cfg, repos.AdminUsers, audit),
		Audit:     audit,
		Razorpay:  razorpay,
		Uploads:    NewUploadService(cfg, repos.Media, audit, storage),
		HeroSlides: NewHeroSlideService(repos.HeroSlides, audit),
		Content:    NewContentService(db, audit),
		Gallery:    NewGalleryService(repos.Gallery, audit),
		Restaurant: NewRestaurantService(repos.Restaurant, audit),
		Suites:     NewSuiteService(repos.Suites, audit),
		Halls:      NewHallService(repos.Halls, audit),
		SEO:        NewSEOService(repos.SEO, audit),
		Offers:     NewOfferService(repos.Offers, audit),
	}
}
