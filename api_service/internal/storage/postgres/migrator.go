package postgres

import (
	"ai-assistant-api/internal/config"
	"ai-assistant-api/internal/model"

	"gorm.io/gorm"
)

type migrator struct {
	db  *gorm.DB
	cfg *config.Config
}

func newMigrator(cfg *config.Config, db *gorm.DB) *migrator {
	return &migrator{
		db:  db,
		cfg: cfg,
	}
}

func (m *migrator) migrate() error {
	return m.db.AutoMigrate(
		&model.Task{},
	)
}
