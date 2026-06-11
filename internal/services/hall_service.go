package services

import (
	"context"
	"time"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type HallService struct {
	repo  *repositories.HallRepository
	audit *AuditService
}

func NewHallService(repo *repositories.HallRepository, audit *AuditService) *HallService {
	return &HallService{repo: repo, audit: audit}
}

type HallPackageInput struct {
	Title         string         `json:"title" validate:"required"`
	Description   string         `json:"description"`
	ImageURL      string         `json:"image_url" validate:"required"`
	Capacity      int            `json:"capacity"`
	StartingPrice int64          `json:"starting_price"`
	SuitableFor   datatypes.JSON `json:"suitable_for"`
	Features      datatypes.JSON `json:"features"`
	IsActive      bool           `json:"is_active"`
	SortOrder     int            `json:"sort_order"`
}

type HallEnquiryInput struct {
	Name           string    `json:"name" validate:"required"`
	Phone          string    `json:"phone" validate:"required"`
	Email          string    `json:"email"`
	EventType      string    `json:"event_type"`
	ExpectedGuests int       `json:"expected_guests"`
	PreferredDate  time.Time `json:"preferred_date"`
	Message        string    `json:"message"`
	Source         string    `json:"source"`
}

func (s *HallService) ListPackagesPublic(ctx context.Context) ([]models.HallPackage, error) {
	return s.repo.ListPackagesPublic(ctx)
}

func (s *HallService) ListPackagesAdmin(ctx context.Context) ([]models.HallPackage, error) {
	return s.repo.ListPackagesAdmin(ctx)
}

func (s *HallService) CreatePackage(ctx context.Context, input HallPackageInput, adminID uuid.UUID, ip string) (*models.HallPackage, error) {
	pkg := &models.HallPackage{
		Title:         input.Title,
		Description:   input.Description,
		ImageURL:      input.ImageURL,
		Capacity:      input.Capacity,
		StartingPrice: input.StartingPrice,
		SuitableFor:   input.SuitableFor,
		Features:      input.Features,
		IsActive:      input.IsActive,
		SortOrder:     input.SortOrder,
	}
	if err := s.repo.CreatePackage(ctx, pkg); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "hall.package_create", "halls", map[string]interface{}{"id": pkg.ID, "title": pkg.Title}, ip)
	return pkg, nil
}

func (s *HallService) UpdatePackage(ctx context.Context, id uuid.UUID, input HallPackageInput, adminID uuid.UUID, ip string) (*models.HallPackage, error) {
	pkg, err := s.repo.FindPackageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	pkg.Title = input.Title
	pkg.Description = input.Description
	pkg.ImageURL = input.ImageURL
	pkg.Capacity = input.Capacity
	pkg.StartingPrice = input.StartingPrice
	pkg.SuitableFor = input.SuitableFor
	pkg.Features = input.Features
	pkg.IsActive = input.IsActive
	pkg.SortOrder = input.SortOrder

	if err := s.repo.SavePackage(ctx, pkg); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "hall.package_update", "halls", map[string]interface{}{"id": pkg.ID}, ip)
	return pkg, nil
}

func (s *HallService) DeletePackage(ctx context.Context, id uuid.UUID, adminID uuid.UUID, ip string) error {
	pkg, err := s.repo.FindPackageByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.DeletePackage(ctx, pkg); err != nil {
		return err
	}
	s.audit.Log(ctx, &adminID, "hall.package_delete", "halls", map[string]interface{}{"id": id, "title": pkg.Title}, ip)
	return nil
}

func (s *HallService) CreateEnquiry(ctx context.Context, input HallEnquiryInput) (*models.HallEnquiry, error) {
	enquiry := &models.HallEnquiry{
		Name:           input.Name,
		Phone:          input.Phone,
		Email:          input.Email,
		EventType:      input.EventType,
		ExpectedGuests: input.ExpectedGuests,
		PreferredDate:  input.PreferredDate,
		Message:        input.Message,
		Source:         input.Source,
		Status:         "new",
	}
	if err := s.repo.CreateEnquiry(ctx, enquiry); err != nil {
		return nil, err
	}
	return enquiry, nil
}

func (s *HallService) ListEnquiries(ctx context.Context, page, limit int) ([]models.HallEnquiry, int64, error) {
	return s.repo.ListEnquiries(ctx, page, limit)
}

func (s *HallService) UpdateEnquiryStatus(ctx context.Context, id uuid.UUID, status string, adminID uuid.UUID, ip string) (*models.HallEnquiry, error) {
	enquiry, err := s.repo.FindEnquiryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	enquiry.Status = status
	if err := s.repo.SaveEnquiry(ctx, enquiry); err != nil {
		return nil, err
	}
	s.audit.Log(ctx, &adminID, "hall.enquiry_status_update", "halls", map[string]interface{}{"id": enquiry.ID, "status": status}, ip)
	return enquiry, nil
}
