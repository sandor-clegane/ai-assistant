package sql

import (
	"ai-assistant-api/internal/model"

	"gorm.io/gorm"
)

type migrator struct {
	db *gorm.DB
}

func newMigrator(db *gorm.DB) *migrator {
	return &migrator{
		db: db,
	}
}

func (m *migrator) migrate() error {
	return m.db.AutoMigrate(
		&model.Task{},
	)
}
