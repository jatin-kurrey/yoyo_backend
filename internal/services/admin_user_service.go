package services

import (
	"context"
	"errors"
	"strings"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminUserService struct {
	cfg   *config.Config
	repo  *repositories.AdminUserRepository
	audit *AuditService
}

type AdminUserInput struct {
	Name     string           `json:"name" validate:"required,min=2,max=120"`
	Email    string           `json:"email" validate:"required,email,max=180"`
	Password string           `json:"password" validate:"omitempty,min=8,max=128"`
	Role     models.AdminRole `json:"role" validate:"required,oneof=super_admin admin moderator staff"`
	IsActive *bool            `json:"is_active"`
}

func NewAdminUserService(cfg *config.Config, repo *repositories.AdminUserRepository, audit *AuditService) *AdminUserService {
	return &AdminUserService{cfg: cfg, repo: repo, audit: audit}
}

func (s *AdminUserService) List(ctx context.Context, search string, page int, limit int) ([]models.AdminUser, int64, error) {
	return s.repo.List(ctx, search, page, limit)
}

func (s *AdminUserService) Create(ctx context.Context, input AdminUserInput, adminID *uuid.UUID, ip string) (*models.AdminUser, error) {
	hash, err := utils.HashPassword(input.Password, s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}
	active := true
	if input.IsActive != nil {
		active = *input.IsActive
	}
	user := &models.AdminUser{
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: hash,
		Role:         input.Role,
		IsActive:     active,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "create", "admin_users", map[string]interface{}{"admin_user_id": user.ID}, ip)
	return user, nil
}

func (s *AdminUserService) Find(ctx context.Context, id uuid.UUID) (*models.AdminUser, error) {
	user, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return user, err
}

func (s *AdminUserService) Update(ctx context.Context, id uuid.UUID, input AdminUserInput, adminID *uuid.UUID, ip string) (*models.AdminUser, error) {
	user, err := s.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	active := user.IsActive
	if input.IsActive != nil {
		active = *input.IsActive
	}
	user.Name = strings.TrimSpace(input.Name)
	user.Email = strings.ToLower(strings.TrimSpace(input.Email))
	user.Role = input.Role
	user.IsActive = active
	if input.Password != "" {
		hash, err := utils.HashPassword(input.Password, s.cfg.BcryptCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, adminID, "update", "admin_users", map[string]interface{}{"admin_user_id": user.ID}, ip)
	return user, nil
}

func (s *AdminUserService) Delete(ctx context.Context, id uuid.UUID, adminID *uuid.UUID, ip string) error {
	user, err := s.Find(ctx, id)
	if err != nil {
		return err
	}
	if user.Role == models.RoleSuperAdmin && user.IsActive {
		count, err := s.repo.CountSuperAdmins(ctx)
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrOnlySuperAdmin
		}
	}
	user.IsActive = false
	if err := s.repo.Save(ctx, user); err != nil {
		return err
	}
	s.audit.Log(ctx, adminID, "deactivate", "admin_users", map[string]interface{}{"admin_user_id": user.ID}, ip)
	return nil
}
