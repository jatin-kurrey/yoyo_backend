package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	cfg   *config.Config
	users *repositories.AdminUserRepository
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResult struct {
	Token string            `json:"token"`
	User  *models.AdminUser `json:"user"`
}

func NewAuthService(cfg *config.Config, users *repositories.AdminUserRepository) *AuthService {
	return &AuthService{cfg: cfg, users: users}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !user.IsActive {
		return nil, ErrInactiveAccount
	}
	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	user.LastLogin = &now
	if err := s.users.Save(ctx, user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateAdminToken(user.ID, user.Email, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTAccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: user}, nil
}

func (s *AuthService) EnsureSuperAdmin(ctx context.Context) error {
	if s.cfg.AdminEmail == "" || s.cfg.AdminPassword == "" {
		return nil
	}

	existing, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(s.cfg.AdminEmail)))
	if err == nil && existing.ID.String() != "" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := utils.HashPassword(s.cfg.AdminPassword, s.cfg.BcryptCost)
	if err != nil {
		return err
	}

	return s.users.Create(ctx, &models.AdminUser{
		Name:         s.cfg.AdminName,
		Email:        strings.ToLower(strings.TrimSpace(s.cfg.AdminEmail)),
		PasswordHash: hash,
		Role:         models.RoleSuperAdmin,
		IsActive:     true,
	})
}
