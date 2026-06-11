package services

import (
	"context"
	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentService struct {
	db    *gorm.DB
	audit *AuditService
}

func NewContentService(db *gorm.DB, audit *AuditService) *ContentService {
	return &ContentService{db: db, audit: audit}
}

func (s *ContentService) FindBySlug(ctx context.Context, slug string) (*models.ContentPage, error) {
	var page models.ContentPage
	if err := s.db.WithContext(ctx).Where("slug = ? AND is_published = ?", slug, true).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &page, nil
}

func (s *ContentService) List(ctx context.Context) ([]models.ContentPage, error) {
	var pages []models.ContentPage
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

func (s *ContentService) AdminFindBySlug(ctx context.Context, slug string) (*models.ContentPage, error) {
	var page models.ContentPage
	if err := s.db.WithContext(ctx).Where("slug = ?", slug).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &page, nil
}

type UpdateContentInput struct {
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	IsPublished *bool   `json:"is_published"`
}

func (s *ContentService) Update(ctx context.Context, adminID uuid.UUID, slug string, input UpdateContentInput, ip string) (*models.ContentPage, error) {
	var page models.ContentPage
	if err := s.db.WithContext(ctx).Where("slug = ?", slug).First(&page).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if input.Title != nil {
		page.Title = *input.Title
	}
	if input.Content != nil {
		page.Content = *input.Content
	}
	if input.IsPublished != nil {
		page.IsPublished = *input.IsPublished
	}

	if err := s.db.WithContext(ctx).Save(&page).Error; err != nil {
		return nil, err
	}

	s.audit.Log(ctx, &adminID, "update", "content", map[string]interface{}{
		"slug":  slug,
		"title": page.Title,
	}, ip)

	return &page, nil
}
