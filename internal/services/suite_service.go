package services

import (
	"context"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SuiteService struct {
	repo  *repositories.SuiteRepository
	audit *AuditService
}

func NewSuiteService(repo *repositories.SuiteRepository, audit *AuditService) *SuiteService {
	return &SuiteService{repo: repo, audit: audit}
}

type SuiteInput struct {
	Title         string         `json:"title" validate:"required"`
	Slug          string         `json:"slug" validate:"required"`
	Description   string         `json:"description"`
	ImageURL      string         `json:"image_url" validate:"required"`
	Gallery       datatypes.JSON `json:"gallery"`
	PricePerNight int64          `json:"price_per_night" validate:"required"`
	MaxGuests     int            `json:"max_guests"`
	Amenities     datatypes.JSON `json:"amenities"`
	IsActive      bool           `json:"is_active"`
	SortOrder     int            `json:"sort_order"`
}

func (s *SuiteService) ListPublic(ctx context.Context) ([]models.SuiteRoom, error) {
	return s.repo.ListPublic(ctx)
}

func (s *SuiteService) ListAdmin(ctx context.Context) ([]models.SuiteRoom, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *SuiteService) FindBySlug(ctx context.Context, slug string) (*models.SuiteRoom, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *SuiteService) Create(ctx context.Context, input SuiteInput, adminID uuid.UUID, ip string) (*models.SuiteRoom, error) {
	suite := &models.SuiteRoom{
		Title:         input.Title,
		Slug:          input.Slug,
		Description:   input.Description,
		ImageURL:      input.ImageURL,
		Gallery:       input.Gallery,
		PricePerNight: input.PricePerNight,
		MaxGuests:     input.MaxGuests,
		Amenities:     input.Amenities,
		IsActive:      input.IsActive,
		SortOrder:     input.SortOrder,
	}
	if err := s.repo.Create(ctx, suite); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "suite.create", "suites", map[string]interface{}{"id": suite.ID, "slug": suite.Slug}, ip)
	return suite, nil
}

func (s *SuiteService) Update(ctx context.Context, id uuid.UUID, input SuiteInput, adminID uuid.UUID, ip string) (*models.SuiteRoom, error) {
	suite, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	suite.Title = input.Title
	suite.Slug = input.Slug
	suite.Description = input.Description
	suite.ImageURL = input.ImageURL
	suite.Gallery = input.Gallery
	suite.PricePerNight = input.PricePerNight
	suite.MaxGuests = input.MaxGuests
	suite.Amenities = input.Amenities
	suite.IsActive = input.IsActive
	suite.SortOrder = input.SortOrder

	if err := s.repo.Save(ctx, suite); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "suite.update", "suites", map[string]interface{}{"id": suite.ID}, ip)
	return suite, nil
}

func (s *SuiteService) Delete(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	suite, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, suite); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "suite.delete", "suites", map[string]interface{}{"id": id, "slug": suite.Slug}, ip)
	return nil
}
