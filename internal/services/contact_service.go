package services

import (
	"context"
	"errors"

	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactService struct {
	repo  *repositories.ContactMessageRepository
	audit *AuditService
}

type ContactInput struct {
	Name    string `json:"name" validate:"required,min=2,max=140"`
	Email   string `json:"email" validate:"required,email,max=180"`
	Phone   string `json:"phone" validate:"omitempty,max=30"`
	Subject string `json:"subject" validate:"omitempty,max=180"`
	Message string `json:"message" validate:"required,min=5,max=5000"`
}

func NewContactService(repo *repositories.ContactMessageRepository, audit *AuditService) *ContactService {
	return &ContactService{repo: repo, audit: audit}
}

func (s *ContactService) Create(ctx context.Context, input ContactInput) (*models.ContactMessage, error) {
	message := &models.ContactMessage{
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Subject: input.Subject,
		Message: input.Message,
		Status:  models.MessageNew,
	}
	return message, s.repo.Create(ctx, message)
}

func (s *ContactService) ListAdmin(ctx context.Context, filter repositories.ContactMessageFilter, page int, limit int) ([]models.ContactMessage, int64, error) {
	return s.repo.ListAdmin(ctx, filter, page, limit)
}

func (s *ContactService) FindAdminByID(ctx context.Context, id uuid.UUID) (*models.ContactMessage, error) {
	message, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return message, err
}

func (s *ContactService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.MessageStatus, adminID *uuid.UUID, ip string) (*models.ContactMessage, error) {
	message, err := s.FindAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	message.Status = status
	if err := s.repo.Save(ctx, message); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "update_status", "messages", map[string]interface{}{"message_id": message.ID, "status": status}, ip)
	return message, nil
}

func (s *ContactService) Delete(ctx context.Context, id uuid.UUID, adminID *uuid.UUID, ip string) error {
	message, err := s.FindAdminByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, message); err != nil {
		return err
	}
	s.audit.Log(ctx, adminID, "delete", "messages", map[string]interface{}{"message_id": message.ID}, ip)
	return nil
}
