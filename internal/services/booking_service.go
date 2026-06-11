package services

import (
	"context"
	"errors"
	"time"

	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService struct {
	db       *gorm.DB
	tickets  *repositories.TicketRepository
	bookings *repositories.BookingRepository
	razorpay *RazorpayService
	audit    *AuditService
}

type CreateOrderInput struct {
	CustomerName  string    `json:"customer_name" validate:"required,min=2,max=140"`
	CustomerEmail string    `json:"customer_email" validate:"required,email,max=180"`
	CustomerPhone string    `json:"customer_phone" validate:"required,min=8,max=30"`
	TicketID      uuid.UUID `json:"ticket_id" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required,gte=1"`
	VisitDate     time.Time `json:"visit_date" validate:"required"`
}

type CreateOrderResult struct {
	Booking       *models.Booking `json:"booking"`
	RazorpayOrder *RazorpayOrder  `json:"razorpay_order"`
	KeyID         string          `json:"key_id"`
}

type VerifyPaymentInput struct {
	RazorpayOrderID   string `json:"razorpay_order_id" validate:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" validate:"required"`
	RazorpaySignature string `json:"razorpay_signature" validate:"required"`
}

func NewBookingService(db *gorm.DB, tickets *repositories.TicketRepository, bookings *repositories.BookingRepository, razorpay *RazorpayService, audit *AuditService) *BookingService {
	return &BookingService{db: db, tickets: tickets, bookings: bookings, razorpay: razorpay, audit: audit}
}

func (s *BookingService) CreateOrder(ctx context.Context, input CreateOrderInput) (*CreateOrderResult, error) {
	ticket, err := s.tickets.FindByID(ctx, input.TicketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if !ticket.IsActive {
		return nil, ErrNotFound
	}
	if ticket.Stock < input.Quantity {
		return nil, ErrInsufficientStock
	}

	bookingCode := utils.NewBookingCode()
	amount := ticket.Price * int64(input.Quantity) * 100
	order, err := s.razorpay.CreateOrder(ctx, amount, bookingCode)
	if err != nil {
		return nil, err
	}

	booking := &models.Booking{
		BookingID:       bookingCode,
		CustomerName:    input.CustomerName,
		CustomerEmail:   input.CustomerEmail,
		CustomerPhone:   input.CustomerPhone,
		TicketID:        ticket.ID,
		Quantity:        input.Quantity,
		Amount:          amount,
		PaymentStatus:   models.PaymentPending,
		RazorpayOrderID: order.ID,
		VisitDate:       input.VisitDate,
		Status:          models.BookingPending,
	}
	if err := s.bookings.Create(ctx, booking); err != nil {
		return nil, err
	}
	booking.Ticket = *ticket

	return &CreateOrderResult{Booking: booking, RazorpayOrder: order, KeyID: s.razorpay.cfg.RazorpayKeyID}, nil
}

func (s *BookingService) VerifyPayment(ctx context.Context, input VerifyPaymentInput) (*models.Booking, error) {
	var verified *models.Booking
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		booking, err := s.bookings.FindByRazorpayOrderIDForUpdate(ctx, tx, input.RazorpayOrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if booking.PaymentStatus == models.PaymentPaid {
			verified = booking
			return nil
		}

		if !s.razorpay.VerifyPaymentSignature(input.RazorpayOrderID, input.RazorpayPaymentID, input.RazorpaySignature) {
			booking.PaymentStatus = models.PaymentFailed
			booking.RazorpayPaymentID = input.RazorpayPaymentID
			booking.RazorpaySignature = input.RazorpaySignature
			_ = s.bookings.SaveWithTx(ctx, tx, booking)
			return ErrInvalidSignature
		}

		ticket, err := s.tickets.FindByIDForUpdate(ctx, tx, booking.TicketID)
		if err != nil {
			return err
		}
		if ticket.Stock < booking.Quantity {
			return ErrInsufficientStock
		}

		ticket.Stock -= booking.Quantity
		ticket.SoldCount += booking.Quantity
		if err := tx.WithContext(ctx).Save(ticket).Error; err != nil {
			return err
		}

		booking.PaymentStatus = models.PaymentPaid
		booking.Status = models.BookingConfirmed
		booking.RazorpayPaymentID = input.RazorpayPaymentID
		booking.RazorpaySignature = input.RazorpaySignature
		if err := s.bookings.SaveWithTx(ctx, tx, booking); err != nil {
			return err
		}
		booking.Ticket = *ticket
		verified = booking
		return nil
	})
	return verified, err
}

func (s *BookingService) ListAdmin(ctx context.Context, filter repositories.BookingFilter, page int, limit int) ([]models.Booking, int64, error) {
	return s.bookings.ListAdmin(ctx, filter, page, limit)
}

func (s *BookingService) FindAdminByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	booking, err := s.bookings.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return booking, err
}

func (s *BookingService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookingStatus, adminID *uuid.UUID, ip string) (*models.Booking, error) {
	booking, err := s.FindAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	booking.Status = status
	if status == models.BookingRefunded {
		booking.PaymentStatus = models.PaymentRefunded
	}
	if err := s.db.WithContext(ctx).Save(booking).Error; err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "update_status", "bookings", map[string]interface{}{"booking_id": booking.ID, "status": status}, ip)
	return booking, nil
}

func (s *BookingService) MarkPaymentFailedByOrder(ctx context.Context, orderID string, paymentID string) error {
	if orderID == "" {
		return nil
	}
	updates := map[string]interface{}{
		"payment_status": models.PaymentFailed,
	}
	if paymentID != "" {
		updates["razorpay_payment_id"] = paymentID
	}
	return s.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("razorpay_order_id = ? AND payment_status <> ?", orderID, models.PaymentPaid).
		Updates(updates).Error
}
