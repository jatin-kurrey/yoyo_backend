package services

import (
	"context"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SEOService struct {
	repo  *repositories.SEORepository
	audit *AuditService
}

func NewSEOService(repo *repositories.SEORepository, audit *AuditService) *SEOService {
	return &SEOService{repo: repo, audit: audit}
}

type SEOInput struct {
	PageSlug        string         `json:"page_slug" validate:"required"`
	MetaTitle       string         `json:"meta_title"`
	MetaDescription string         `json:"meta_description"`
	CanonicalURL    string         `json:"canonical_url"`
	OGTitle         string         `json:"og_title"`
	OGDescription   string         `json:"og_description"`
	OGImage         string         `json:"og_image"`
	RobotsIndex     bool           `json:"robots_index"`
	RobotsFollow    bool           `json:"robots_follow"`
	SchemaJSON      datatypes.JSON `json:"schema_json"`
}

func (s *SEOService) GetPublic(ctx context.Context, slug string) (*models.SEOPage, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *SEOService) ListAdmin(ctx context.Context) ([]models.SEOPage, error) {
	return s.repo.List(ctx)
}

func (s *SEOService) Save(ctx context.Context, input SEOInput, adminID uuid.UUID, ip string) (*models.SEOPage, error) {
	page, err := s.repo.FindBySlug(ctx, input.PageSlug)
	if err != nil {
		// Create new if not found
		page = &models.SEOPage{PageSlug: input.PageSlug}
	}

	page.MetaTitle = input.MetaTitle
	page.MetaDescription = input.MetaDescription
	page.CanonicalURL = input.CanonicalURL
	page.OGTitle = input.OGTitle
	page.OGDescription = input.OGDescription
	page.OGImage = input.OGImage
	page.RobotsIndex = input.RobotsIndex
	page.RobotsFollow = input.RobotsFollow
	page.SchemaJSON = input.SchemaJSON

	if err := s.repo.Save(ctx, page); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "seo.update", "seo", map[string]interface{}{"slug": page.PageSlug}, ip)
	return page, nil
}

func (s *SEOService) GetByID(ctx context.Context, id uuid.UUID) (*models.SEOPage, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *SEOService) Delete(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	page, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, page); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "seo.delete", "seo", map[string]interface{}{"slug": page.PageSlug}, ip)
	return nil
}
