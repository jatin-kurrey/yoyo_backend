package services

import (
	"context"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
)

type AttractionService struct {
	repo  *repositories.AttractionRepository
	audit *AuditService
}

func NewAttractionService(repo *repositories.AttractionRepository, audit *AuditService) *AttractionService {
	return &AttractionService{repo: repo, audit: audit}
}

type AttractionInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IconName    string `json:"icon_name"`
	Tag         string `json:"tag"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

func (s *AttractionService) ListPublic(ctx context.Context) ([]models.Attraction, error) {
	return s.repo.ListPublic(ctx)
}

func (s *AttractionService) ListAdmin(ctx context.Context) ([]models.Attraction, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *AttractionService) Create(ctx context.Context, input AttractionInput, adminID uuid.UUID, ip string) (*models.Attraction, error) {
	item := &models.Attraction{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		IconName:    input.IconName,
		Tag:         input.Tag,
		SortOrder:   input.SortOrder,
		IsActive:    input.IsActive,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "attraction.create", "attractions", map[string]interface{}{"id": item.ID, "title": item.Title}, ip)
	return item, nil
}

func (s *AttractionService) Update(ctx context.Context, id uuid.UUID, input AttractionInput, adminID uuid.UUID, ip string) (*models.Attraction, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Title = input.Title
	item.Description = input.Description
	item.ImageURL = input.ImageURL
	item.IconName = input.IconName
	item.Tag = input.Tag
	item.SortOrder = input.SortOrder
	item.IsActive = input.IsActive

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "attraction.update", "attractions", map[string]interface{}{"id": item.ID}, ip)
	return item, nil
}

func (s *AttractionService) Delete(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, item); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "attraction.delete", "attractions", map[string]interface{}{"id": id, "title": item.Title}, ip)
	return nil
}
