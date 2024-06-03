package sql

import (
	"ai-assistant-api/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(cfg *config.Config) (*Storage, error) {
	db, err := gorm.Open(sqlite.Open(cfg.StoragePath), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	err = newMigrator(db).migrate()
	if err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}, nil
}
