package postgres

import (
	"ai-assistant-api/internal/config"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func New(db *gorm.DB, getter *trmgorm.CtxGetter) *Storage {
	return &Storage{
		db:     db,
		getter: getter,
	}
}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.StoragePath), &gorm.Config{
		TranslateError: true,
	})
}

func MustNewPostgresDB(cfg *config.Config) *gorm.DB {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	err = newMigrator(cfg, db).migrate()
	if err != nil {
		panic(err)
	}

	return db
}
