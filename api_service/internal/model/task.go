package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Status      TaskStatus `gorm:"column:task_status"`
	Instruction string     `gorm:"column:instruction"`
	Code        string     `gorm:"column:code"`
	Response    string     `gorm:"column:response"`

	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		t.ID = id
	}

	return nil
}
