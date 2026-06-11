package database

import (
	"time"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Warn
	if cfg.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if cfg.AutoMigrate && cfg.AppEnv != "production" {
		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
			return nil, err
		}
		if err := db.AutoMigrate(
			&models.AdminUser{},
			&models.Ticket{},
			&models.Booking{},
			&models.ContactMessage{},
			&models.SiteSetting{},
			&models.AuditLog{},
			&models.HeroSlide{},
			&models.ContentPage{},
			&models.MediaAsset{},
			&models.SEOPage{},
			&models.GalleryItem{},
			&models.RestaurantItem{},
			&models.SuiteRoom{},
			&models.HallPackage{},
			&models.HallEnquiry{},
			&models.Offer{},
		); err != nil {
			return nil, err
		}
	}

	return db, nil
}
