// Enchanted-Garden/internal/repository/repository.go
package repository

import (
	"Enchanted-Garden/internal/model"
	"context"
)

// BranchRepository описывает работу с ветвями (департаментами) в базе данных
type BranchRepository interface {
	Create(ctx context.Context, branch *model.Branch) error
	GetByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error)
	Update(ctx context.Context, branch *model.Branch) error
	DeleteCascade(ctx context.Context, id uint) error
	DeleteReassign(ctx context.Context, id uint, reassignToID uint) error
	FindDuplicate(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error)
	CheckCycle(ctx context.Context, id uint, newParentID uint) (bool, error)
}

// FloraRepository описывает работу с флорой (сотрудниками) в базе данных
type FloraRepository interface {
	Create(ctx context.Context, flora *model.Flora) error
}
