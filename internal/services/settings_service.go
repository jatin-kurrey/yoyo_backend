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
	SiteName        string                 `json:"site_name" validate:"required,min=2,max=140"`
	LogoURL         string                 `json:"logo_url" validate:"omitempty,max=1000"`
	ContactEmail    string                 `json:"contact_email" validate:"omitempty,email,max=180"`
	PhoneNumbers    []string               `json:"phone_numbers"`
	Address         string                 `json:"address" validate:"omitempty,max=1000"`
	SocialLinks     map[string]string      `json:"social_links"`
	SEOTitle        string                 `json:"seo_title" validate:"omitempty,max=180"`
	SEODescription  string                 `json:"seo_description" validate:"omitempty,max=500"`
	RazorpayEnabled bool                   `json:"razorpay_enabled"`
	MaintenanceMode bool                   `json:"maintenance_mode"`
	FeatureToggles  map[string]interface{} `json:"feature_toggles"`
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
	setting := &models.SiteSetting{
		SiteName:         "YOYO FUN N FOODS",
		ContactEmail:     "hello@yoyofun.com",
		PhoneNumbers:     datatypes.JSON(phones),
		Address:          "YOYO FUN N FOODS, Madhya Pradesh, India",
		SocialLinks:      datatypes.JSON(social),
		MetaTitle:        "YOYO FUN N FOODS - Water Park Booking",
		MetaDescription:  "Book tickets for YOYO FUN N FOODS and enjoy a safe, fun-filled park experience.",
		RazorpayEnabled:  true,
		MaintenanceMode:  false,
		FeatureToggles:   datatypes.JSON(toggles),
		HomepageSections: datatypes.JSON(homeSections),
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

	if err := s.repo.Save(ctx, setting); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "update", "settings", nil, ip)
	return setting, nil
}

func PublicSettings(setting *models.SiteSetting) map[string]interface{} {
	return map[string]interface{}{
		"site_name":        setting.SiteName,
		"logo_url":         setting.LogoURL,
		"contact_email":    setting.ContactEmail,
		"phone_numbers":    setting.PhoneNumbers,
		"address":          setting.Address,
		"social_links":     setting.SocialLinks,
		"seo_title":        setting.MetaTitle,
		"seo_description":  setting.MetaDescription,
		"razorpay_enabled": setting.RazorpayEnabled,
		"maintenance_mode": setting.MaintenanceMode,
		"feature_toggles":  setting.FeatureToggles,
	}
}
