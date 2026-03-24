package app

import (
	"parser/internal/config"
	"parser/internal/repository"

	"github.com/rs/zerolog"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	bookRepository *repository.Repository
}
