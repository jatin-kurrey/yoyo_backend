package seeds

import (
	"context"
	"encoding/json"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Run(ctx context.Context, cfg *config.Config, db *gorm.DB, svc *services.Services) error {
	if err := svc.Auth.EnsureSuperAdmin(ctx); err != nil {
		return err
	}
	if _, err := svc.Settings.Get(ctx); err != nil {
		return err
	}
	return seedTickets(ctx, cfg, db)
}

func seedTickets(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	_ = cfg
	defaults := []struct {
		Title         string
		Description   string
		Price         int64
		OriginalPrice *int64
		Category      string
		Stock         int
		SortOrder     int
		Features      []string
	}{
		{
			Title:       "Standard Pass",
			Description: "Single entry to the park with all-day access to core attractions.",
			Price:       499,
			Category:    "general",
			Stock:       100,
			SortOrder:   1,
			Features:    []string{"All Day Entry", "Locker Access", "Safety Gear Included"},
		},
		{
			Title:       "VIP Pass",
			Description: "Skip the line experience with a complimentary drink for a smoother day out.",
			Price:       999,
			Category:    "vip",
			Stock:       50,
			SortOrder:   2,
			Features:    []string{"Priority Entry", "Free Drink", "Locker Access", "Safety Gear Included"},
		},
		{
			Title:       "Family Bundle",
			Description: "Entry for 4 people, designed for families and small groups.",
			Price:       2999,
			Category:    "family",
			Stock:       30,
			SortOrder:   3,
			Features:    []string{"Entry for 4", "Family Check-in", "Locker Access", "Safety Gear Included"},
		},
	}

	for _, item := range defaults {
		slug := utils.Slugify(item.Title)
		var existing models.Ticket
		err := db.WithContext(ctx).Where("slug = ?", slug).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		features, err := json.Marshal(item.Features)
		if err != nil {
			return err
		}
		ticket := models.Ticket{
			Title:         item.Title,
			Slug:          slug,
			Description:   item.Description,
			Price:         item.Price,
			OriginalPrice: item.OriginalPrice,
			Category:      item.Category,
			Features:      datatypes.JSON(features),
			Validity:      "Valid for selected visit date",
			Stock:         item.Stock,
			IsActive:      true,
			SortOrder:     item.SortOrder,
		}
		if err := db.WithContext(ctx).Create(&ticket).Error; err != nil {
			return err
		}
	}
	return nil
}
