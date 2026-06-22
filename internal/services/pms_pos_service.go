package services

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PMSPOSService struct {
	db *gorm.DB
}

func NewPMSPOSService(db *gorm.DB) *PMSPOSService {
	return &PMSPOSService{db: db}
}

func (s *PMSPOSService) ListTables(ctx context.Context) ([]models.POSTable, error) {
	var tables []models.POSTable
	err := s.db.WithContext(ctx).Order("table_number ASC").Find(&tables).Error
	return tables, err
}

func (s *PMSPOSService) OccupyTable(ctx context.Context, tableID uuid.UUID, guestName string) (*models.POSTable, error) {
	var table models.POSTable
	if err := s.db.WithContext(ctx).First(&table, "id = ?", tableID).Error; err != nil {
		return nil, ErrNotFound
	}
	table.Status = models.POSTableOccupied
	table.GuestName = guestName
	table.KOTCount = 0
	table.CurrentOrderValue = 0
	if err := s.db.WithContext(ctx).Save(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

type AddKOTInput struct {
	TableID    uuid.UUID `json:"table_id" validate:"required"`
	MenuItemID *uuid.UUID `json:"menu_item_id"`
	ItemName   string    `json:"item_name" validate:"required"`
	Quantity   int       `json:"quantity"`
	Price      int64     `json:"price" validate:"required"`
	Notes      string    `json:"notes"`
}

func (s *PMSPOSService) AddKOT(ctx context.Context, input AddKOTInput) (*models.POSOrder, error) {
	var table models.POSTable
	if err := s.db.WithContext(ctx).First(&table, "id = ?", input.TableID).Error; err != nil {
		return nil, ErrNotFound
	}
	if input.Quantity < 1 {
		input.Quantity = 1
	}

	nextKOT := table.KOTCount + 1

	order := &models.POSOrder{
		TableID:    input.TableID,
		MenuItemID: input.MenuItemID,
		ItemName:   input.ItemName,
		Quantity:   input.Quantity,
		Price:      input.Price,
		Notes:      input.Notes,
		KOTNumber:  nextKOT,
		Status:     models.POSOrderOpen,
		CreatedAt:  time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}

	totalAdd := input.Price * int64(input.Quantity)
	s.db.WithContext(ctx).Model(&models.POSTable{}).Where("id = ?", input.TableID).
		Updates(map[string]interface{}{
			"kot_count":          nextKOT,
			"current_order_value": gorm.Expr("current_order_value + ?", totalAdd),
		})

	return order, nil
}

func (s *PMSPOSService) GenerateBill(ctx context.Context, tableID uuid.UUID) (*models.POSTable, error) {
	var table models.POSTable
	if err := s.db.WithContext(ctx).First(&table, "id = ?", tableID).Error; err != nil {
		return nil, ErrNotFound
	}
	table.Status = models.POSTableBilled
	if err := s.db.WithContext(ctx).Save(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (s *PMSPOSService) VacateTable(ctx context.Context, tableID uuid.UUID) (*models.POSTable, error) {
	var table models.POSTable
	if err := s.db.WithContext(ctx).First(&table, "id = ?", tableID).Error; err != nil {
		return nil, ErrNotFound
	}
	table.Status = models.POSTableVacant
	table.GuestName = ""
	table.KOTCount = 0
	table.CurrentOrderValue = 0
	if err := s.db.WithContext(ctx).Save(&table).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (s *PMSPOSService) MoveToRoomFolio(ctx context.Context, tableID uuid.UUID, bookingID uuid.UUID) error {
	var table models.POSTable
	if err := s.db.WithContext(ctx).First(&table, "id = ?", tableID).Error; err != nil {
		return ErrNotFound
	}

	entry := &models.FolioEntry{
		BookingID:   bookingID,
		Type:        models.FolioFood,
		Description: "POS Table " + string(rune('0'+table.TableNumber)) + " - Restaurant",
		Amount:      table.CurrentOrderValue,
		Quantity:    1,
		PostedAt:    time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(entry).Error; err != nil {
		return err
	}

	s.db.WithContext(ctx).Model(&models.PMSBooking{}).Where("id = ?", bookingID).
		Update("balance_amount", gorm.Expr("balance_amount + ?", table.CurrentOrderValue))

	table.Status = models.POSTableVacant
	table.GuestName = ""
	table.KOTCount = 0
	table.CurrentOrderValue = 0
	return s.db.WithContext(ctx).Save(&table).Error
}

func (s *PMSPOSService) GetKOTs(ctx context.Context, tableID uuid.UUID) ([]models.POSOrder, error) {
	var orders []models.POSOrder
	err := s.db.WithContext(ctx).Where("table_id = ?", tableID).Order("created_at ASC").Find(&orders).Error
	return orders, err
}
