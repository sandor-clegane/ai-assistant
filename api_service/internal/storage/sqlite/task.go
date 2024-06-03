package sql

import (
	"ai-assistant-api/internal/model"
	"ai-assistant-api/internal/storage"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Storage) SaveTask(ctx context.Context, task model.Task) (*model.Task, error) {
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"task_status", "response"}),
	}).Create(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Storage) GetTask(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	var task model.Task

	err := s.db.Model(&model.Task{}).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (s *Storage) TaskListByStatus(ctx context.Context, status model.TaskStatus) ([]model.Task, error) {
	var taskList []model.Task

	err := s.db.Where("task_status = ?", status).Find(&taskList).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return taskList, nil
}
