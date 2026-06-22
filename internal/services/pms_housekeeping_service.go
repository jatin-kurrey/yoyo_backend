package services

import (
	"context"
	"time"

	"yoyo-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PMSHKService struct {
db *gorm.DB
}

func NewPMSHKService(db *gorm.DB) *PMSHKService {
	return &PMSHKService{db: db}
}

type CreateHKTaskInput struct {
	RoomID    uuid.UUID        `json:"room_id" validate:"required"`
	StaffName string           `json:"staff_name"`
	TaskType  models.HKTaskType `json:"task_type" validate:"required"`
	Notes     string           `json:"notes"`
}

func (s *PMSHKService) ListTasks(ctx context.Context, status string) ([]models.HKTask, error) {
	var tasks []models.HKTask
	query := s.db.WithContext(ctx).Preload("Room").Preload("Room.Category")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (s *PMSHKService) CreateTask(ctx context.Context, input CreateHKTaskInput) (*models.HKTask, error) {
	task := &models.HKTask{
		RoomID:    input.RoomID,
		StaffName: input.StaffName,
		TaskType:  input.TaskType,
		Status:    models.HKPending,
		Notes:     input.Notes,
	}
	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (s *PMSHKService) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status models.HKTaskStatus) (*models.HKTask, error) {
	var task models.HKTask
	if err := s.db.WithContext(ctx).First(&task, "id = ?", taskID).Error; err != nil {
		return nil, ErrNotFound
	}
	task.Status = status
	now := time.Now()
	switch status {
	case models.HKInProgress:
		task.AssignedAt = &now
	case models.HKCompleted:
		task.CompletedAt = &now
	}
	if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
		return nil, err
	}

	if status == models.HKCompleted && task.TaskType == models.HKClean {
		s.db.WithContext(ctx).Model(&models.PMSRoom{}).Where("id = ?", task.RoomID).
			Update("clean_status", models.CleanClean)
	}

	return &task, nil
}

func (s *PMSHKService) SetRoomCleanStatus(ctx context.Context, roomID uuid.UUID, clean models.CleanStatus) (*models.PMSRoom, error) {
	var room models.PMSRoom
	if err := s.db.WithContext(ctx).First(&room, "id = ?", roomID).Error; err != nil {
		return nil, ErrNotFound
	}
	room.CleanStatus = clean
	if err := s.db.WithContext(ctx).Save(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *PMSHKService) SetRoomOOO(ctx context.Context, roomID uuid.UUID, reason string) (*models.PMSRoom, error) {
	var room models.PMSRoom
	if err := s.db.WithContext(ctx).First(&room, "id = ?", roomID).Error; err != nil {
		return nil, ErrNotFound
	}
	room.Status = models.RoomOOO
	room.OOOReason = reason
	if err := s.db.WithContext(ctx).Save(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *PMSHKService) SetRoomAvailable(ctx context.Context, roomID uuid.UUID) (*models.PMSRoom, error) {
	var room models.PMSRoom
	if err := s.db.WithContext(ctx).First(&room, "id = ?", roomID).Error; err != nil {
		return nil, ErrNotFound
	}
	room.Status = models.RoomAvailable
	room.OOOReason = ""
	if err := s.db.WithContext(ctx).Save(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
