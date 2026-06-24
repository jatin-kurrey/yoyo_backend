package services

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PMSAccountsService struct {
	db *gorm.DB
}

func NewPMSAccountsService(db *gorm.DB) *PMSAccountsService {
	return &PMSAccountsService{db: db}
}

func (s *PMSAccountsService) ListTransactions(ctx context.Context, txType, status string) ([]models.PMSTransaction, error) {
	var list []models.PMSTransaction
	query := s.db.WithContext(ctx).Order("created_at DESC")
	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&list).Error
	return list, err
}

type CreateTransactionInput struct {
	Date        string `json:"date" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=income expense"`
	Category    string `json:"category" validate:"required"`
	Description string `json:"description"`
	Amount      int64  `json:"amount" validate:"required,min=1"`
	Method      string `json:"method" validate:"required"`
	Status      string `json:"status" validate:"omitempty,oneof=completed pending"`
	GuestName   string `json:"guest_name"`
}

func (s *PMSAccountsService) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*models.PMSTransaction, error) {
	if input.Status == "" {
		input.Status = "completed"
	}
	tx := &models.PMSTransaction{
		Date:        input.Date,
		Type:        input.Type,
		Category:    input.Category,
		Description: input.Description,
		Amount:      input.Amount,
		Method:      input.Method,
		Status:      input.Status,
		GuestName:   input.GuestName,
		CreatedAt:   time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(tx).Error; err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *PMSAccountsService) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&models.PMSTransaction{}, "id = ?", id).Error
}

// ── Settings ──

func (s *PMSAccountsService) GetSettings(ctx context.Context) (map[string]string, error) {
	var settings []models.PMSSetting
	if err := s.db.WithContext(ctx).Find(&settings).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

func (s *PMSAccountsService) UpsertSetting(ctx context.Context, key, value string) error {
	var existing models.PMSSetting
	err := s.db.WithContext(ctx).Where("key = ?", key).First(&existing).Error
	if err == nil {
		return s.db.WithContext(ctx).Model(&existing).Update("value", value).Error
	}
	return s.db.WithContext(ctx).Create(&models.PMSSetting{Key: key, Value: value}).Error
}

// ── Rate Overrides ──

type RateOverrideInput struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	Date       string    `json:"date" validate:"required"`
	Plan       string    `json:"plan" validate:"required,oneof=ep cp ap"`
	Rate       int64     `json:"rate"`
	StopSell   *bool     `json:"stop_sell"`
}

func (s *PMSAccountsService) SetRateOverride(ctx context.Context, input RateOverrideInput) (*models.PMSRateOverride, error) {
	var existing models.PMSRateOverride
	err := s.db.WithContext(ctx).
		Where("category_id = ? AND date = ? AND plan = ?", input.CategoryID, input.Date, input.Plan).
		First(&existing).Error

	if err == nil {
		if input.Rate > 0 {
			existing.Rate = input.Rate
		}
		if input.StopSell != nil {
			existing.StopSell = *input.StopSell
		}
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}

	override := &models.PMSRateOverride{
		CategoryID: input.CategoryID,
		Date:       input.Date,
		Plan:       input.Plan,
		Rate:       input.Rate,
		StopSell:   input.StopSell != nil && *input.StopSell,
	}
	if err := s.db.WithContext(ctx).Create(override).Error; err != nil {
		return nil, err
	}
	return override, nil
}

func (s *PMSAccountsService) ClearRateOverride(ctx context.Context, categoryID uuid.UUID, date, plan string) error {
	return s.db.WithContext(ctx).
		Where("category_id = ? AND date = ? AND plan = ?", categoryID, date, plan).
		Delete(&models.PMSRateOverride{}).Error
}

func (s *PMSAccountsService) ListRateOverrides(ctx context.Context, categoryID *uuid.UUID) ([]models.PMSRateOverride, error) {
	var list []models.PMSRateOverride
	query := s.db.WithContext(ctx)
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	err := query.Order("date ASC").Find(&list).Error
	return list, err
}
