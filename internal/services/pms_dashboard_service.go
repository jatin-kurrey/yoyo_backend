package services

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"gorm.io/gorm"
)

type PMSDashboardService struct {
	db *gorm.DB
}

type PMSDashboardStats struct {
	TotalRooms     int64   `json:"total_rooms"`
	OccupiedRooms  int64   `json:"occupied_rooms"`
	VacantRooms    int64   `json:"vacant_rooms"`
	OOORooms       int64   `json:"ooo_rooms"`
	OccupancyRate  float64 `json:"occupancy_rate"`
	TotalRevenue   int64   `json:"total_revenue"`
	TodayRevenue   int64   `json:"today_revenue"`
	TodayArrivals  int64   `json:"today_arrivals"`
	TodayDepartures int64  `json:"today_departures"`
	CleanRooms     int64   `json:"clean_rooms"`
	DirtyRooms     int64   `json:"dirty_rooms"`
	InHouse        int64   `json:"in_house"`
	PendingBalance int64   `json:"pending_balance"`
}

func NewPMSDashboardService(db *gorm.DB) *PMSDashboardService {
	return &PMSDashboardService{db: db}
}

func (s *PMSDashboardService) GetStats(ctx context.Context) (*PMSDashboardStats, error) {
	stats := &PMSDashboardStats{}
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Count(&stats.TotalRooms)
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("status = ?", models.RoomOccupied).Count(&stats.OccupiedRooms)
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("status = ?", models.RoomOOO).Count(&stats.OOORooms)
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("status = ?", models.RoomAvailable).Count(&stats.VacantRooms)
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("clean_status = ?", models.CleanClean).Count(&stats.CleanRooms)
	s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("clean_status = ?", models.CleanDirty).Count(&stats.DirtyRooms)

	s.db.WithContext(ctx).Model(&models.PMSBooking{}).Where("status = ?", models.PMSBookingCheckedIn).Count(&stats.InHouse)

	s.db.WithContext(ctx).Model(&models.PMSBooking{}).
		Where("check_in >= ? AND check_in < ?", today, tomorrow).
		Count(&stats.TodayArrivals)

	s.db.WithContext(ctx).Model(&models.PMSBooking{}).
		Where("check_out >= ? AND check_out < ?", today, tomorrow).
		Count(&stats.TodayDepartures)

	s.db.WithContext(ctx).Model(&models.PMSPayment{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TotalRevenue)

	s.db.WithContext(ctx).Model(&models.PMSPayment{}).
		Where("received_at >= ? AND received_at < ?", today, tomorrow).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TodayRevenue)

	s.db.WithContext(ctx).Model(&models.PMSBooking{}).
		Select("COALESCE(SUM(balance_amount), 0)").
		Scan(&stats.PendingBalance)

	if stats.TotalRooms > 0 {
		stats.OccupancyRate = float64(stats.OccupiedRooms) / float64(stats.TotalRooms) * 100
	}

	return stats, nil
}
