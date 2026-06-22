package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PMSBookingService struct {
	db *gorm.DB
}

func NewPMSBookingService(db *gorm.DB) *PMSBookingService {
	return &PMSBookingService{db: db}
}

type CreatePMSBookingInput struct {
	RoomID       uuid.UUID              `json:"room_id" validate:"required"`
	GuestName    string                 `json:"guest_name" validate:"required"`
	GuestPhone   string                 `json:"guest_phone" validate:"required"`
	GuestEmail   string                 `json:"guest_email"`
	Adults       int                    `json:"adults"`
	Children     int                    `json:"children"`
	Plan         models.BookingPlan     `json:"plan"`
	Source       models.BookingSource   `json:"source"`
	CheckIn      time.Time              `json:"check_in" validate:"required"`
	CheckOut     time.Time              `json:"check_out" validate:"required"`
	RatePerNight int64                  `json:"rate_per_night"`
	Discount     int64                  `json:"discount"`
	Status       models.BookingStatusPMS `json:"status"`
}

func (s *PMSBookingService) Create(ctx context.Context, input CreatePMSBookingInput, adminID *uuid.UUID) (*models.PMSBooking, error) {
	if input.Status == "" {
		input.Status = models.PMSBookingCheckedIn
	}
	if input.Source == "" {
		input.Source = models.SourceWalkIn
	}
	if input.Plan == "" {
		input.Plan = models.PlanEP
	}

	var room models.PMSRoom
	if err := s.db.WithContext(ctx).First(&room, "id = ?", input.RoomID).Error; err != nil {
		return nil, ErrNotFound
	}
	if input.RatePerNight == 0 {
		var cat models.PMSRoomCategory
		if err := s.db.WithContext(ctx).First(&cat, "id = ?", room.CategoryID).Error; err == nil {
			input.RatePerNight = cat.BasePrice
		}
	}

	nights := int(math.Ceil(input.CheckOut.Sub(input.CheckIn).Hours() / 24))
	if nights < 1 {
		nights = 1
	}
	total := input.RatePerNight * int64(nights)
	tax := int64(math.Round(float64(total) * 0.12))
	grandTotal := total + tax - input.Discount

	ref := fmt.Sprintf("YOYO-%s-%d", time.Now().Format("20060102"), time.Now().UnixMilli()%10000)

	booking := &models.PMSBooking{
		BookingRef:   ref,
		RoomID:       input.RoomID,
		GuestName:    input.GuestName,
		GuestPhone:   input.GuestPhone,
		GuestEmail:   input.GuestEmail,
		Adults:       input.Adults,
		Children:     input.Children,
		Plan:         input.Plan,
		Source:       input.Source,
		CheckIn:      input.CheckIn,
		CheckOut:     input.CheckOut,
		RatePerNight: input.RatePerNight,
		TotalAmount:  grandTotal,
		Discount:     input.Discount,
		Tax:          tax,
		BalanceAmount: grandTotal,
		Status:       input.Status,
		CreatedByID:  adminID,
	}

	if err := s.db.WithContext(ctx).Create(booking).Error; err != nil {
		return nil, err
	}

		if input.Status == models.PMSBookingCheckedIn {
		s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("id = ?", input.RoomID).
			Update("status", models.RoomOccupied)
	}

	return booking, nil
}

func (s *PMSBookingService) List(ctx context.Context, search string, status string, page, limit int) ([]models.PMSBooking, int64, error) {
	var bookings []models.PMSBooking
	query := s.db.WithContext(ctx).Model(&models.PMSBooking{}).Preload("Room").Preload("Room.Category")

	if search != "" {
		term := "%" + search + "%"
		query = query.Where("guest_name ILIKE ? OR guest_phone ILIKE ? OR booking_ref ILIKE ?", term, term, term)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&bookings).Error
	return bookings, total, err
}

func (s *PMSBookingService) GetByID(ctx context.Context, id uuid.UUID) (*models.PMSBooking, error) {
	var booking models.PMSBooking
	err := s.db.WithContext(ctx).Preload("Room").Preload("Room.Category").
		First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return &booking, nil
}

func (s *PMSBookingService) CheckIn(ctx context.Context, id uuid.UUID) (*models.PMSBooking, error) {
	var booking models.PMSBooking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		return nil, ErrNotFound
	}
	booking.Status = models.PMSBookingCheckedIn
	if err := s.db.WithContext(ctx).Save(&booking).Error; err != nil {
		return nil, err
	}
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("id = ?", booking.RoomID).
		Update("status", models.RoomOccupied)
	return &booking, nil
}

func (s *PMSBookingService) CheckOut(ctx context.Context, id uuid.UUID) (*models.PMSBooking, error) {
	var booking models.PMSBooking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		return nil, ErrNotFound
	}
	booking.Status = models.PMSBookingCheckedOut
	if err := s.db.WithContext(ctx).Save(&booking).Error; err != nil {
		return nil, err
	}
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("id = ?", booking.RoomID).
		Update("status", models.RoomAvailable)
	return &booking, nil
}

func (s *PMSBookingService) Cancel(ctx context.Context, id uuid.UUID) (*models.PMSBooking, error) {
	var booking models.PMSBooking
	if err := s.db.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		return nil, ErrNotFound
	}
	booking.Status = models.PMSBookingCancelled
	if err := s.db.WithContext(ctx).Save(&booking).Error; err != nil {
		return nil, err
	}
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("id = ?", booking.RoomID).
		Update("status", models.RoomAvailable)
	return &booking, nil
}
