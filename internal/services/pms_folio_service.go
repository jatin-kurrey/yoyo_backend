package services

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PMSBookingsService struct {
	db *gorm.DB
}

func NewPMSFolioService(db *gorm.DB) *PMSBookingsService {
	return &PMSBookingsService{db: db}
}

type AddFolioEntryInput struct {
	BookingID   uuid.UUID            `json:"booking_id" validate:"required"`
	Type        models.FolioEntryType `json:"type" validate:"required"`
	Description string               `json:"description"`
	Amount      int64                `json:"amount" validate:"required"`
	Quantity    int                  `json:"quantity"`
}

type AddPaymentInput struct {
	BookingID uuid.UUID           `json:"booking_id" validate:"required"`
	Mode      models.PaymentMode  `json:"mode" validate:"required"`
	Amount    int64               `json:"amount" validate:"required"`
	Type      models.PaymentType  `json:"type"`
	Reference string              `json:"reference"`
}

func (s *PMSBookingsService) GetFolio(ctx context.Context, bookingID uuid.UUID) ([]models.FolioEntry, []models.PMSPayment, error) {
	var entries []models.FolioEntry
	if err := s.db.WithContext(ctx).Where("booking_id = ?", bookingID).Order("posted_at ASC").Find(&entries).Error; err != nil {
		return nil, nil, err
	}
	var payments []models.PMSPayment
	if err := s.db.WithContext(ctx).Where("booking_id = ?", bookingID).Order("received_at ASC").Find(&payments).Error; err != nil {
		return nil, nil, err
	}
	return entries, payments, nil
}

func (s *PMSBookingsService) AddFolioEntry(ctx context.Context, input AddFolioEntryInput, adminID *uuid.UUID) (*models.FolioEntry, error) {
	var booking models.PMSBooking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", input.BookingID).Error; err != nil {
		return nil, ErrNotFound
	}
	if input.Quantity < 1 {
		input.Quantity = 1
	}
	entry := &models.FolioEntry{
		BookingID:   input.BookingID,
		Type:        input.Type,
		Description: input.Description,
		Amount:      input.Amount,
		Quantity:    input.Quantity,
		PostedAt:    time.Now(),
		PostedByID:  adminID,
	}
	if err := s.db.WithContext(ctx).Create(entry).Error; err != nil {
		return nil, err
	}
	totalCharge := input.Amount * int64(input.Quantity)
	s.db.WithContext(ctx).Model(&models.PMSBooking{}).Where("id = ?", input.BookingID).
		Update("balance_amount", gorm.Expr("balance_amount + ?", totalCharge))
	return entry, nil
}

func (s *PMSBookingsService) AddPayment(ctx context.Context, input AddPaymentInput, adminID *uuid.UUID) (*models.PMSPayment, error) {
	var booking models.PMSBooking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", input.BookingID).Error; err != nil {
		return nil, ErrNotFound
	}
	if input.Type == "" {
		input.Type = models.PayAdvance
	}
	payment := &models.PMSPayment{
		BookingID:    input.BookingID,
		Mode:         input.Mode,
		Amount:       input.Amount,
		Type:         input.Type,
		Reference:    input.Reference,
		ReceivedAt:   time.Now(),
		ReceivedByID: adminID,
	}
	if err := s.db.WithContext(ctx).Create(payment).Error; err != nil {
		return nil, err
	}
	s.db.WithContext(ctx).Model(&models.PMSBooking{}).Where("id = ?", input.BookingID).
		Updates(map[string]interface{}{
			"paid_amount":   gorm.Expr("paid_amount + ?", input.Amount),
			"balance_amount": gorm.Expr("balance_amount - ?", input.Amount),
		})
	return payment, nil
}
