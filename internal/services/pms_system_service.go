package services

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"gorm.io/gorm"
)

type PMSSystemService struct {
	db *gorm.DB
}

func NewPMSSystemService(db *gorm.DB) *PMSSystemService {
	return &PMSSystemService{db: db}
}

type BackupData struct {
	ExportedAt   string                    `json:"exported_at"`
	Version      string                    `json:"version"`
	Categories   []models.PMSRoomCategory  `json:"categories"`
	Rooms        []models.PMSRoom          `json:"rooms"`
	Bookings     []models.PMSBooking       `json:"bookings"`
	FolioEntries []models.FolioEntry       `json:"folio_entries"`
	Payments     []models.PMSPayment       `json:"payments"`
	POSTables    []models.POSTable         `json:"pos_tables"`
	POSOrders    []models.POSOrder         `json:"pos_orders"`
	HKTasks      []models.HKTask           `json:"hk_tasks"`
	Transactions []models.PMSTransaction   `json:"transactions"`
	Settings     []models.PMSSetting       `json:"settings"`
	RateOverrides []models.PMSRateOverride `json:"rate_overrides"`
}

func (s *PMSSystemService) Export(ctx context.Context) (*BackupData, error) {
	data := &BackupData{
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Version:    "2.0",
	}

	var err error
	data.Categories, err = s.listCategories(ctx)
	if err != nil {
		return nil, err
	}
	data.Rooms, err = s.listRooms(ctx)
	if err != nil {
		return nil, err
	}
	data.Bookings, err = s.listBookings(ctx)
	if err != nil {
		return nil, err
	}
	data.FolioEntries, err = s.listFolioEntries(ctx)
	if err != nil {
		return nil, err
	}
	data.Payments, err = s.listPayments(ctx)
	if err != nil {
		return nil, err
	}
	data.POSTables, err = s.listPOSTables(ctx)
	if err != nil {
		return nil, err
	}
	data.POSOrders, err = s.listPOSOrders(ctx)
	if err != nil {
		return nil, err
	}
	data.HKTasks, err = s.listHKTasks(ctx)
	if err != nil {
		return nil, err
	}
	data.Transactions, err = s.listTransactions(ctx)
	if err != nil {
		return nil, err
	}
	data.Settings, err = s.listSettings(ctx)
	if err != nil {
		return nil, err
	}
	data.RateOverrides, err = s.listRateOverrides(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PMSSystemService) Import(ctx context.Context, data *BackupData) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(data.Categories) > 0 {
			for i := range data.Categories {
				data.Categories[i].ID = data.Categories[i].ID
				data.Categories[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.Categories[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.Rooms) > 0 {
			for i := range data.Rooms {
				data.Rooms[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.Rooms[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.Bookings) > 0 {
			for i := range data.Bookings {
				data.Bookings[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.Bookings[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.FolioEntries) > 0 {
			for i := range data.FolioEntries {
				if err := tx.Create(&data.FolioEntries[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.Payments) > 0 {
			for i := range data.Payments {
				if err := tx.Create(&data.Payments[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.POSTables) > 0 {
			for i := range data.POSTables {
				data.POSTables[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.POSTables[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.POSOrders) > 0 {
			for i := range data.POSOrders {
				if err := tx.Create(&data.POSOrders[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.HKTasks) > 0 {
			for i := range data.HKTasks {
				data.HKTasks[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.HKTasks[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.Transactions) > 0 {
			for i := range data.Transactions {
				data.Transactions[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.Transactions[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.Settings) > 0 {
			for i := range data.Settings {
				if err := tx.Create(&data.Settings[i]).Error; err != nil {
					return err
				}
			}
		}
		if len(data.RateOverrides) > 0 {
			for i := range data.RateOverrides {
				data.RateOverrides[i].DeletedAt = gorm.DeletedAt{}
				if err := tx.Create(&data.RateOverrides[i]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *PMSSystemService) ResetAll(ctx context.Context) error {
	tables := []interface{}{
		&models.PMSRateOverride{},
		&models.PMSTransaction{},
		&models.PMSSetting{},
		&models.HKTask{},
		&models.POSOrder{},
		&models.POSTable{},
		&models.PMSPayment{},
		&models.FolioEntry{},
		&models.PMSBooking{},
		&models.PMSRoom{},
		&models.PMSRoomCategory{},
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, table := range tables {
			if err := tx.Where("1 = 1").Delete(table).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *PMSSystemService) GetStats(ctx context.Context) (map[string]int64, error) {
	stats := map[string]int64{}
	tx := s.db.WithContext(ctx)

	var count int64
	tx.Model(&models.PMSRoomCategory{}).Count(&count)
	stats["categories"] = count

	tx.Model(&models.PMSRoom{}).Count(&count)
	stats["rooms"] = count

	tx.Model(&models.PMSBooking{}).Count(&count)
	stats["bookings"] = count

	tx.Model(&models.PMSTransaction{}).Count(&count)
	stats["transactions"] = count

	tx.Model(&models.PMSRateOverride{}).Count(&count)
	stats["rate_overrides"] = count

	tx.Model(&models.POSTable{}).Count(&count)
	stats["pos_tables"] = count

	tx.Model(&models.HKTask{}).Count(&count)
	stats["hk_tasks"] = count

	tx.Model(&models.PMSSetting{}).Count(&count)
	stats["settings"] = count

	return stats, nil
}

// Internal list helpers
func (s *PMSSystemService) listCategories(ctx context.Context) ([]models.PMSRoomCategory, error) {
	var list []models.PMSRoomCategory
	err := s.db.WithContext(ctx).Unscoped().Order("created_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listRooms(ctx context.Context) ([]models.PMSRoom, error) {
	var list []models.PMSRoom
	err := s.db.WithContext(ctx).Unscoped().Preload("Category").Order("room_number ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listBookings(ctx context.Context) ([]models.PMSBooking, error) {
	var list []models.PMSBooking
	err := s.db.WithContext(ctx).Unscoped().Preload("Room").Order("created_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listFolioEntries(ctx context.Context) ([]models.FolioEntry, error) {
	var list []models.FolioEntry
	err := s.db.WithContext(ctx).Order("posted_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listPayments(ctx context.Context) ([]models.PMSPayment, error) {
	var list []models.PMSPayment
	err := s.db.WithContext(ctx).Order("received_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listPOSTables(ctx context.Context) ([]models.POSTable, error) {
	var list []models.POSTable
	err := s.db.WithContext(ctx).Unscoped().Order("table_number ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listPOSOrders(ctx context.Context) ([]models.POSOrder, error) {
	var list []models.POSOrder
	err := s.db.WithContext(ctx).Order("created_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listHKTasks(ctx context.Context) ([]models.HKTask, error) {
	var list []models.HKTask
	err := s.db.WithContext(ctx).Unscoped().Order("created_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listTransactions(ctx context.Context) ([]models.PMSTransaction, error) {
	var list []models.PMSTransaction
	err := s.db.WithContext(ctx).Unscoped().Order("created_at ASC").Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listSettings(ctx context.Context) ([]models.PMSSetting, error) {
	var list []models.PMSSetting
	err := s.db.WithContext(ctx).Find(&list).Error
	return list, err
}
func (s *PMSSystemService) listRateOverrides(ctx context.Context) ([]models.PMSRateOverride, error) {
	var list []models.PMSRateOverride
	err := s.db.WithContext(ctx).Unscoped().Order("date ASC").Find(&list).Error
	return list, err
}
