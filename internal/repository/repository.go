package repository

import (
	"context"
	"errors"
	"fmt"

	"parser/internal/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) GetBookById(ctx context.Context, bookID uint64) (model.ParseBook, error) {
	var book model.ParseBook
	res := r.db.WithContext(ctx).
		Model(&model.ParseBook{}).
		Where("id = ? AND is_deleted = ?", bookID, false).
		First(&book)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.ParseBook{}, fmt.Errorf("book not found")
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to get book id")
	}
	return book, nil
}
