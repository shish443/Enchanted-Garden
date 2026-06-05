// Enchanted-Garden/internal/repository/postgres/branch.go
package postgres

import (
	"Enchanted-Garden/internal/model"
	"context"

	"gorm.io/gorm"
)

type BranchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(ctx context.Context, branch *model.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *BranchRepository) Update(ctx context.Context, branch *model.Branch) error {
	return r.db.WithContext(ctx).Save(branch).Error
}

func (r *BranchRepository) DeleteCascade(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Branch{}, id).Error
}

func (r *BranchRepository) DeleteReassign(ctx context.Context, id uint, reassignToID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var currentBranch model.Branch
		if err := tx.First(&currentBranch, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Flora{}).Where("branch_id = ?", id).Update("branch_id", reassignToID).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Branch{}).Where("parent_id = ?", id).Update("parent_id", currentBranch.ParentID).Error; err != nil {
			return err
		}

		return tx.Delete(&model.Branch{}, id).Error
	})
}

func (r *BranchRepository) GetByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error) {
	var branch model.Branch
	query := r.db.WithContext(ctx)

	if sortBy != "created_at" && sortBy != "full_name" {
		sortBy = "created_at"
	}
	sortQuery := sortBy + " ASC"

	if includeEmployees {
		query = query.Preload("Flora", func(db *gorm.DB) *gorm.DB {
			return db.Order(sortQuery)
		})
	}

	currentPath := "Children"
	for i := 1; i <= depth; i++ {
		query = query.Preload(currentPath)
		if includeEmployees {
			query = query.Preload(currentPath+".Flora", func(db *gorm.DB) *gorm.DB {
				return db.Order(sortQuery)
			})
		}
		currentPath += ".Children"
	}

	err := query.First(&branch, id).Error
	if err != nil {
		return nil, err
	}

	return &branch, nil
}
func (r *BranchRepository) CheckCycle(ctx context.Context, id uint, newParentID uint) (bool, error) {
	if id == newParentID {
		return true, nil
	}

	var count int64
	query := `
		WITH RECURSIVE tree AS (
			SELECT id FROM branches WHERE id = ?
			UNION ALL
			SELECT b.id FROM branches b
			JOIN tree t ON b.parent_id = t.id
		)
		SELECT COUNT(*) FROM tree WHERE id = ?
	`
	if err := r.db.WithContext(ctx).Raw(query, newParentID, id).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *BranchRepository) FindDuplicate(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Branch{}).Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
