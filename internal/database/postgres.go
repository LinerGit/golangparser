package database

import (
	"log"
	"parser/internal/config"
	"parser/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}
	logger := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
	if err != nil {
		logger.Error().Msgf("failed to conn to db")
		return nil, err
	}
	logger.Info().Msg("db connected")
	return db, err
}
