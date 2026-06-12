package services

import (
	"context"
	"encoding/json"
	"errors"

	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SettingsService struct {
	repo  *repositories.SiteSettingRepository
	audit *AuditService
}

type SettingsInput struct {
	SiteName            string                 `json:"site_name" validate:"required,min=2,max=140"`
	LogoURL             string                 `json:"logo_url" validate:"omitempty,max=1000"`
	ContactEmail        string                 `json:"contact_email" validate:"omitempty,email,max=180"`
	PhoneNumbers        []string               `json:"phone_numbers"`
	Address             string                 `json:"address" validate:"omitempty,max=1000"`
	SocialLinks         map[string]string      `json:"social_links"`
	SEOTitle            string                 `json:"seo_title" validate:"omitempty,max=180"`
	SEODescription      string                 `json:"seo_description" validate:"omitempty,max=500"`
	RazorpayEnabled     bool                   `json:"razorpay_enabled"`
	MaintenanceMode     bool                   `json:"maintenance_mode"`
	FeatureToggles      map[string]interface{} `json:"feature_toggles"`
	AdminSidebarToggles map[string]interface{} `json:"admin_sidebar_toggles"`
	AboutHeadline       string                 `json:"about_headline"`
	AboutDescription    string                 `json:"about_description"`
	AboutVideoURL       string                 `json:"about_video_url"`
	AboutImage1URL      string                 `json:"about_image_1_url"`
	AboutImage2URL      string                 `json:"about_image_2_url"`
	AboutBullets        interface{}            `json:"about_bullets"`
	TrustBullets        interface{}            `json:"trust_bullets"`
}

func NewSettingsService(repo *repositories.SiteSettingRepository, audit *AuditService) *SettingsService {
	return &SettingsService{repo: repo, audit: audit}
}

func (s *SettingsService) Get(ctx context.Context) (*models.SiteSetting, error) {
	setting, err := s.repo.First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.CreateDefault(ctx)
	}
	return setting, err
}

func (s *SettingsService) CreateDefault(ctx context.Context) (*models.SiteSetting, error) {
	phones, _ := json.Marshal([]string{"+91 9752586956", "+91 9589986956"})
	social, _ := json.Marshal(map[string]string{
		"facebook":  "",
		"instagram": "",
	})
	toggles, _ := json.Marshal(map[string]interface{}{
		"onlineBooking": true,
		"contactForm":   true,
	})
	homeSections, _ := json.Marshal(map[string]bool{
		"hero":       true,
		"tickets":    true,
		"restaurant": true,
		"suites":     true,
		"halls":      true,
		"gallery":    true,
		"offers":     true,
		"contact":    true,
	})
	sidebarToggles, _ := json.Marshal(map[string]bool{
		"Dashboard":      true,
		"Hero Landing":   true,
		"Tickets":        true,
		"Bookings":       true,
		"Messages":       true,
		"Gallery":        true,
		"Attractions":    true,
		"Restaurant":     true,
		"Suites & Rooms": true,
		"Events & Halls": true,
		"Promotions":     true,
		"Content Mgmt":   true,
		"SEO Manager":    true,
		"Settings":       true,
		"Users":          true,
		"Audit Logs":     true,
	})
	aboutB, _ := json.Marshal([]map[string]string{
		{"icon": "Smile", "title": "Family-Friendly Fun", "desc": "Safe rides, clean facilities, and activities designed for all ages to enjoy together."},
		{"icon": "Lightbulb", "title": "Thrilling Adventures", "desc": "High-speed slides, massive wave pools, and exciting attractions await at every turn."},
	})
	trustB, _ := json.Marshal([]map[string]string{
		{"icon": "ShieldCheck", "title": "Certified Lifeguards", "desc": "Every pool is monitored by trained professionals 24/7."},
		{"icon": "Zap", "title": "Hygiene Protocol", "desc": "Daily water testing and continuous filtration cycles."},
		{"icon": "Heart", "title": "Family Friendly", "desc": "Dedicated lockers and private changing rooms for families."},
		{"icon": "Award", "title": "Award Winning", "desc": "Ranked #1 for safety and cleanliness in the region."},
	})

	setting := &models.SiteSetting{
		SiteName:            "YOYO FUN N FOODS",
		ContactEmail:        "hello@yoyofun.com",
		PhoneNumbers:        datatypes.JSON(phones),
		Address:             "Village Godhi, Tehsil Ahiwara, District Durg, Chhattisgarh 490036",
		SocialLinks:         datatypes.JSON(social),
		MetaTitle:           "YOYO FUN N FOODS - Water Park Booking",
		MetaDescription:     "Book tickets for YOYO FUN N FOODS and enjoy a safe, fun-filled park experience.",
		RazorpayEnabled:     true,
		MaintenanceMode:     false,
		FeatureToggles:      datatypes.JSON(toggles),
		HomepageSections:    datatypes.JSON(homeSections),
		AdminSidebarToggles: datatypes.JSON(sidebarToggles),
		AboutHeadline:       "Central India's Favorite Water Park Destination",
		AboutDescription:    "At YOYO Fun N Foods, every visit is a splash of excitement and joy. Families create unforgettable memories with thrilling rides, delicious food, and endless fun.",
		AboutVideoURL:       "/about-bg.mp4",
		AboutImage1URL:      "https://images.unsplash.com/photo-1629834598512-77a443808b73?q=80&w=687&auto=format&fit=crop",
		AboutImage2URL:      "https://plus.unsplash.com/premium_photo-1661378818245-0fe1239e5bdc?q=80&w=687&auto=format&fit=crop",
		AboutBullets:        datatypes.JSON(aboutB),
		TrustBullets:        datatypes.JSON(trustB),
	}
	return setting, s.repo.Create(ctx, setting)
}

func (s *SettingsService) Update(ctx context.Context, input SettingsInput, adminID *uuid.UUID, ip string) (*models.SiteSetting, error) {
	setting, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}
	phones, err := json.Marshal(input.PhoneNumbers)
	if err != nil {
		return nil, err
	}
	social, err := json.Marshal(input.SocialLinks)
	if err != nil {
		return nil, err
	}
	toggles, err := json.Marshal(input.FeatureToggles)
	if err != nil {
		return nil, err
	}
	sidebarToggles, err := json.Marshal(input.AdminSidebarToggles)
	if err != nil {
		return nil, err
	}
	aboutBullets, err := json.Marshal(input.AboutBullets)
	if err != nil {
		return nil, err
	}
	trustBullets, err := json.Marshal(input.TrustBullets)
	if err != nil {
		return nil, err
	}

	setting.SiteName = input.SiteName
	setting.LogoURL = input.LogoURL
	setting.ContactEmail = input.ContactEmail
	setting.PhoneNumbers = datatypes.JSON(phones)
	setting.Address = input.Address
	setting.SocialLinks = datatypes.JSON(social)
	setting.MetaTitle = input.SEOTitle
	setting.MetaDescription = input.SEODescription
	setting.RazorpayEnabled = input.RazorpayEnabled
	setting.MaintenanceMode = input.MaintenanceMode
	setting.FeatureToggles = datatypes.JSON(toggles)
	setting.AdminSidebarToggles = datatypes.JSON(sidebarToggles)
	setting.AboutHeadline = input.AboutHeadline
	setting.AboutDescription = input.AboutDescription
	setting.AboutVideoURL = input.AboutVideoURL
	setting.AboutImage1URL = input.AboutImage1URL
	setting.AboutImage2URL = input.AboutImage2URL
	setting.AboutBullets = datatypes.JSON(aboutBullets)
	setting.TrustBullets = datatypes.JSON(trustBullets)

	if err := s.repo.Save(ctx, setting); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "update", "settings", nil, ip)
	return setting, nil
}

func PublicSettings(setting *models.SiteSetting) map[string]interface{} {
	return map[string]interface{}{
		"site_name":             setting.SiteName,
		"logo_url":              setting.LogoURL,
		"contact_email":         setting.ContactEmail,
		"phone_numbers":         setting.PhoneNumbers,
		"address":               setting.Address,
		"social_links":          setting.SocialLinks,
		"seo_title":             setting.MetaTitle,
		"seo_description":       setting.MetaDescription,
		"razorpay_enabled":      setting.RazorpayEnabled,
		"maintenance_mode":      setting.MaintenanceMode,
		"feature_toggles":       setting.FeatureToggles,
		"admin_sidebar_toggles": setting.AdminSidebarToggles,
		"about_headline":        setting.AboutHeadline,
		"about_description":     setting.AboutDescription,
		"about_video_url":       setting.AboutVideoURL,
		"about_image_1_url":     setting.AboutImage1URL,
		"about_image_2_url":     setting.AboutImage2URL,
		"about_bullets":         setting.AboutBullets,
		"trust_bullets":         setting.TrustBullets,
	}
}
