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

// создаем простой репозиторий для цветов нашего сада
func NewFloraRepository(db *gorm.DB) *FloraRepository {
	return &FloraRepository{db: db}
}

// Create просто добавляет новую флору (сотрудника) в базу данных
func (r *FloraRepository) Create(ctx context.Context, flora *model.Flora) error {
	return r.db.WithContext(ctx).Create(flora).Error
}
