package services

import (
	"context"
	"time"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
)

type OfferService struct {
	repo  *repositories.OfferRepository
	audit *AuditService
}

func NewOfferService(repo *repositories.OfferRepository, audit *AuditService) *OfferService {
	return &OfferService{repo: repo, audit: audit}
}

type OfferInput struct {
	Title         string     `json:"title" validate:"required"`
	Description   string     `json:"description"`
	Code          string     `json:"code" validate:"required"`
	DiscountType  string     `json:"discount_type" validate:"required"`
	DiscountValue int64      `json:"discount_value" validate:"required"`
	StartsAt      *time.Time `json:"starts_at"`
	EndsAt        *time.Time `json:"ends_at"`
	IsActive      bool       `json:"is_active"`
}

func (s *OfferService) ListActive(ctx context.Context) ([]models.Offer, error) {
	return s.repo.ListActive(ctx)
}

func (s *OfferService) ListAdmin(ctx context.Context) ([]models.Offer, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *OfferService) Create(ctx context.Context, input OfferInput, adminID uuid.UUID, ip string) (*models.Offer, error) {
	offer := &models.Offer{
		Title:         input.Title,
		Description:   input.Description,
		Code:          input.Code,
		DiscountType:  input.DiscountType,
		DiscountValue: input.DiscountValue,
		StartsAt:      input.StartsAt,
		EndsAt:        input.EndsAt,
		IsActive:      input.IsActive,
	}
	if err := s.repo.Create(ctx, offer); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "offer.create", "offers", map[string]interface{}{"id": offer.ID, "code": offer.Code}, ip)
	return offer, nil
}

func (s *OfferService) Update(ctx context.Context, id uuid.UUID, input OfferInput, adminID uuid.UUID, ip string) (*models.Offer, error) {
	offer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	offer.Title = input.Title
	offer.Description = input.Description
	offer.Code = input.Code
	offer.DiscountType = input.DiscountType
	offer.DiscountValue = input.DiscountValue
	offer.StartsAt = input.StartsAt
	offer.EndsAt = input.EndsAt
	offer.IsActive = input.IsActive

	if err := s.repo.Save(ctx, offer); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "offer.update", "offers", map[string]interface{}{"id": offer.ID}, ip)
	return offer, nil
}

func (s *OfferService) Delete(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	offer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, offer); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "offer.delete", "offers", map[string]interface{}{"id": id, "code": offer.Code}, ip)
	return nil
}
