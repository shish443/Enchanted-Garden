// Enchanted-Garden/internal/repository/postgres/flora.go
package postgres

import (
	"Enchanted-Garden/internal/model"
	"context"

	"gorm.io/gorm"
)

type FloraRepository struct {
	db *gorm.DB
}

func NewFloraRepository(db *gorm.DB) *FloraRepository {
	return &FloraRepository{db: db}
}

func (r *FloraRepository) Create(ctx context.Context, flora *model.Flora) error {
	return r.db.WithContext(ctx).Create(flora).Error
}
