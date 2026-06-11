package services

import (
	"context"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
)

type GalleryService struct {
	repo  *repositories.GalleryRepository
	audit *AuditService
}

func NewGalleryService(repo *repositories.GalleryRepository, audit *AuditService) *GalleryService {
	return &GalleryService{repo: repo, audit: audit}
}

type GalleryInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" validate:"required"`
	Category    string `json:"category"`
	AltText     string `json:"alt_text"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

func (s *GalleryService) ListPublic(ctx context.Context) ([]models.GalleryItem, error) {
	return s.repo.ListPublic(ctx)
}

func (s *GalleryService) ListAdmin(ctx context.Context) ([]models.GalleryItem, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *GalleryService) Create(ctx context.Context, input GalleryInput, adminID uuid.UUID, ip string) (*models.GalleryItem, error) {
	item := &models.GalleryItem{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Category:    input.Category,
		AltText:     input.AltText,
		SortOrder:   input.SortOrder,
		IsActive:    input.IsActive,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "gallery.create", "gallery", map[string]interface{}{"id": item.ID, "title": item.Title}, ip)
	return item, nil
}

func (s *GalleryService) Update(ctx context.Context, id uuid.UUID, input GalleryInput, adminID uuid.UUID, ip string) (*models.GalleryItem, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Title = input.Title
	item.Description = input.Description
	item.ImageURL = input.ImageURL
	item.Category = input.Category
	item.AltText = input.AltText
	item.SortOrder = input.SortOrder
	item.IsActive = input.IsActive

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "gallery.update", "gallery", map[string]interface{}{"id": item.ID}, ip)
	return item, nil
}

func (s *GalleryService) Delete(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, item); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "gallery.delete", "gallery", map[string]interface{}{"id": id, "title": item.Title}, ip)
	return nil
}
