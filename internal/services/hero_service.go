package services

import (
	"context"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
)

type HeroSlideService struct {
	repo  *repositories.HeroSlideRepository
	audit *AuditService
}

func NewHeroSlideService(repo *repositories.HeroSlideRepository, audit *AuditService) *HeroSlideService {
	return &HeroSlideService{repo: repo, audit: audit}
}

func (s *HeroSlideService) ListPublic(ctx context.Context) ([]models.HeroSlide, error) {
	return s.repo.ListPublic(ctx)
}

func (s *HeroSlideService) ListAdmin(ctx context.Context) ([]models.HeroSlide, error) {
	return s.repo.ListAdmin(ctx)
}

type HeroSlideInput struct {
	ImageURL    string `json:"image_url" validate:"required,url"`
	Headline    string `json:"headline" validate:"required,min=5,max=255"`
	Subheadline string `json:"subheadline" validate:"max=255"`
	CTAUrl      string `json:"cta_url" validate:"max=255"`
	CTAText     string `json:"cta_text" validate:"max=100"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

func (s *HeroSlideService) Create(ctx context.Context, input HeroSlideInput, adminID *uuid.UUID) (*models.HeroSlide, error) {
	slide := &models.HeroSlide{
		ImageURL:    input.ImageURL,
		Title:       input.Headline,
		Subtitle:    input.Subheadline,
		CTAURL:      input.CTAUrl,
		CTALabel:    input.CTAText,
		IsActive:    input.IsActive,
		SortOrder:   input.SortOrder,
	}

	if err := s.repo.Create(ctx, slide); err != nil {
		return nil, err
	}

	s.audit.Log(ctx, adminID, "create", "hero_slides", map[string]interface{}{
		"slide_id": slide.ID,
		"title":    slide.Title,
	}, "")

	return slide, nil
}

func (s *HeroSlideService) Update(ctx context.Context, id uuid.UUID, input HeroSlideInput, adminID *uuid.UUID) (*models.HeroSlide, error) {
	slide, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	slide.ImageURL = input.ImageURL
	slide.Title = input.Headline
	slide.Subtitle = input.Subheadline
	slide.CTAURL = input.CTAUrl
	slide.CTALabel = input.CTAText
	slide.IsActive = input.IsActive
	slide.SortOrder = input.SortOrder

	if err := s.repo.Save(ctx, slide); err != nil {
		return nil, err
	}

	s.audit.Log(ctx, adminID, "update", "hero_slides", map[string]interface{}{
		"slide_id": slide.ID,
		"title":    slide.Title,
	}, "")

	return slide, nil
}

func (s *HeroSlideService) Delete(ctx context.Context, id uuid.UUID, adminID *uuid.UUID) error {
	slide, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if err := s.repo.Delete(ctx, slide); err != nil {
		return err
	}

	s.audit.Log(ctx, adminID, "delete", "hero_slides", map[string]interface{}{
		"slide_id": id,
		"title":    slide.Title,
	}, "")

	return nil
}
