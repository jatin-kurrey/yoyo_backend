package services

import (
	"context"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PMSPricingService struct {
	db *gorm.DB
}

func NewPMSPricingService(db *gorm.DB) *PMSPricingService {
	return &PMSPricingService{db: db}
}

func (s *PMSPricingService) ListCategoriesWithRooms(ctx context.Context) ([]models.PMSRoomCategory, error) {
	var categories []models.PMSRoomCategory
	err := s.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (s *PMSPricingService) UpdateRates(ctx context.Context, categoryID uuid.UUID, baseRate int64) (*models.PMSRoomCategory, error) {
	var cat models.PMSRoomCategory
	if err := s.db.WithContext(ctx).First(&cat, "id = ?", categoryID).Error; err != nil {
		return nil, ErrNotFound
	}
	cat.BasePrice = baseRate
	if err := s.db.WithContext(ctx).Save(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *PMSPricingService) ListRooms(ctx context.Context) ([]models.PMSRoom, error) {
	var rooms []models.PMSRoom
	err := s.db.WithContext(ctx).Preload("Category").Order("room_number ASC").Find(&rooms).Error
	return rooms, err
}

func (s *PMSPricingService) UpdateRoomStatus(ctx context.Context, roomID uuid.UUID, status models.RoomStatus) (*models.PMSRoom, error) {
	var room models.PMSRoom
	if err := s.db.WithContext(ctx).First(&room, "id = ?", roomID).Error; err != nil {
		return nil, ErrNotFound
	}
	room.Status = status
	if err := s.db.WithContext(ctx).Save(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
