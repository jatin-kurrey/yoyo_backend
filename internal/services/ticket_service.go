package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TicketService struct {
	repo  *repositories.TicketRepository
	audit *AuditService
}

type TicketInput struct {
	Title         string   `json:"title" validate:"required,min=2,max=160"`
	Slug          string   `json:"slug" validate:"omitempty,max=180"`
	Description   string   `json:"description" validate:"omitempty,max=5000"`
	Price         int64    `json:"price" validate:"required,gte=1"`
	OriginalPrice *int64   `json:"original_price" validate:"omitempty,gte=1"`
	Category      string   `json:"category" validate:"omitempty,max=80"`
	Features      []string `json:"features"`
	Validity      string   `json:"validity" validate:"omitempty,max=120"`
	Stock         int      `json:"stock" validate:"gte=0"`
	IsActive      *bool    `json:"is_active"`
	SortOrder     int      `json:"sort_order"`
}

func NewTicketService(repo *repositories.TicketRepository, audit *AuditService) *TicketService {
	return &TicketService{repo: repo, audit: audit}
}

func (s *TicketService) ListPublic(ctx context.Context) ([]models.Ticket, error) {
	return s.repo.ListPublic(ctx)
}

func (s *TicketService) FindPublicBySlug(ctx context.Context, slug string) (*models.Ticket, error) {
	ticket, err := s.repo.FindBySlug(ctx, slug)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return ticket, err
}

func (s *TicketService) ListAdmin(ctx context.Context, filter repositories.TicketFilter, page int, limit int) ([]models.Ticket, int64, error) {
	return s.repo.ListAdmin(ctx, filter, page, limit)
}

func (s *TicketService) FindAdminByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	ticket, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return ticket, err
}

func (s *TicketService) Create(ctx context.Context, input TicketInput, adminID *uuid.UUID, ip string) (*models.Ticket, error) {
	features, err := marshalStringSlice(input.Features)
	if err != nil {
		return nil, err
	}
	active := true
	if input.IsActive != nil {
		active = *input.IsActive
	}
	slug := strings.TrimSpace(input.Slug)
	if slug == "" {
		slug = utils.Slugify(input.Title)
	} else {
		slug = utils.Slugify(slug)
	}

	ticket := &models.Ticket{
		Title:         strings.TrimSpace(input.Title),
		Slug:          slug,
		Description:   strings.TrimSpace(input.Description),
		Price:         input.Price,
		OriginalPrice: input.OriginalPrice,
		Category:      strings.TrimSpace(input.Category),
		Features:      features,
		Validity:      strings.TrimSpace(input.Validity),
		Stock:         input.Stock,
		IsActive:      active,
		SortOrder:     input.SortOrder,
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "create", "tickets", map[string]interface{}{"ticket_id": ticket.ID, "title": ticket.Title}, ip)
	return ticket, nil
}

func (s *TicketService) Update(ctx context.Context, id uuid.UUID, input TicketInput, adminID *uuid.UUID, ip string) (*models.Ticket, error) {
	ticket, err := s.FindAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	features, err := marshalStringSlice(input.Features)
	if err != nil {
		return nil, err
	}
	active := ticket.IsActive
	if input.IsActive != nil {
		active = *input.IsActive
	}

	ticket.Title = strings.TrimSpace(input.Title)
	ticket.Description = strings.TrimSpace(input.Description)
	ticket.Price = input.Price
	ticket.OriginalPrice = input.OriginalPrice
	ticket.Category = strings.TrimSpace(input.Category)
	ticket.Features = features
	ticket.Validity = strings.TrimSpace(input.Validity)
	ticket.Stock = input.Stock
	ticket.IsActive = active
	ticket.SortOrder = input.SortOrder
	if input.Slug != "" {
		ticket.Slug = utils.Slugify(input.Slug)
	}

	if err := s.repo.Save(ctx, ticket); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "update", "tickets", map[string]interface{}{"ticket_id": ticket.ID}, ip)
	return ticket, nil
}

func (s *TicketService) Delete(ctx context.Context, id uuid.UUID, adminID *uuid.UUID, ip string) error {
	ticket, err := s.FindAdminByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, ticket); err != nil {
		return err
	}
	s.audit.Log(ctx, adminID, "delete", "tickets", map[string]interface{}{"ticket_id": ticket.ID}, ip)
	return nil
}

func (s *TicketService) ToggleStatus(ctx context.Context, id uuid.UUID, adminID *uuid.UUID, ip string) (*models.Ticket, error) {
	ticket, err := s.FindAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	ticket.IsActive = !ticket.IsActive
	if err := s.repo.Save(ctx, ticket); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "toggle_status", "tickets", map[string]interface{}{"ticket_id": ticket.ID, "is_active": ticket.IsActive}, ip)
	return ticket, nil
}

func marshalStringSlice(values []string) (datatypes.JSON, error) {
	clean := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			clean = append(clean, value)
		}
	}
	bytes, err := json.Marshal(clean)
	if err != nil {
		return nil, fmt.Errorf("invalid features: %w", err)
	}
	return datatypes.JSON(bytes), nil
}
