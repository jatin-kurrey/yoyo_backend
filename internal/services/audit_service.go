package services

import (
	"context"
	"encoding/json"

	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditService struct {
	repo *repositories.AuditLogRepository
}

func NewAuditService(repo *repositories.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(ctx context.Context, adminID *uuid.UUID, action string, module string, metadata interface{}, ip string) {
	payload := datatypes.JSON([]byte(`{}`))
	if metadata != nil {
		if bytes, err := json.Marshal(metadata); err == nil {
			payload = datatypes.JSON(bytes)
		}
	}
	_ = s.repo.Create(ctx, &models.AuditLog{
		AdminUserID: adminID,
		Action:      action,
		Module:      module,
		Metadata:    payload,
		IPAddress:   ip,
	})
}

func (s *AuditService) List(ctx context.Context, filter repositories.AuditLogFilter, page int, limit int) ([]models.AuditLog, int64, error) {
	return s.repo.List(ctx, filter, page, limit)
}
