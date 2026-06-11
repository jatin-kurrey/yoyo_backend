package services

import (
	"context"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
)

type RestaurantService struct {
	repo  *repositories.RestaurantRepository
	audit *AuditService
}

func NewRestaurantService(repo *repositories.RestaurantRepository, audit *AuditService) *RestaurantService {
	return &RestaurantService{repo: repo, audit: audit}
}

type RestaurantItemInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
	Price       int64  `json:"price"`
	IsFeatured  bool   `json:"is_featured"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

func (s *RestaurantService) ListPublic(ctx context.Context) ([]models.RestaurantItem, error) {
	return s.repo.ListPublic(ctx)
}

func (s *RestaurantService) ListAdmin(ctx context.Context) ([]models.RestaurantItem, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *RestaurantService) CreateItem(ctx context.Context, input RestaurantItemInput, adminID uuid.UUID, ip string) (*models.RestaurantItem, error) {
	item := &models.RestaurantItem{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Category:    input.Category,
		Price:       input.Price,
		IsFeatured:  input.IsFeatured,
		SortOrder:   input.SortOrder,
		IsActive:    input.IsActive,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "restaurant.item_create", "restaurant", map[string]interface{}{"id": item.ID, "title": item.Title}, ip)
	return item, nil
}

func (s *RestaurantService) UpdateItem(ctx context.Context, id uuid.UUID, input RestaurantItemInput, adminID uuid.UUID, ip string) (*models.RestaurantItem, error) {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Title = input.Title
	item.Description = input.Description
	item.ImageURL = input.ImageURL
	item.Category = input.Category
	item.Price = input.Price
	item.IsFeatured = input.IsFeatured
	item.SortOrder = input.SortOrder
	item.IsActive = input.IsActive

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "restaurant.item_update", "restaurant", map[string]interface{}{"id": item.ID}, ip)
	return item, nil
}

func (s *RestaurantService) DeleteItem(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, item); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "restaurant.item_delete", "restaurant", map[string]interface{}{"id": id, "title": item.Title}, ip)
	return nil
}
