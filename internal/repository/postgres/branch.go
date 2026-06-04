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
	// База данных сама удалит всё каскадом, нам достаточно удалить только саму ветку
	return r.db.WithContext(ctx).Delete(&model.Branch{}, id).Error
}

func (r *BranchRepository) DeleteReassign(ctx context.Context, id uint, reassignToID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Находим саму ветку, чтобы узнать, кто её родитель
		var currentBranch model.Branch
		if err := tx.First(&currentBranch, id).Error; err != nil {
			return err
		}

		// Переводим сотрудников этой ветки в новое место
		if err := tx.Model(&model.Flora{}).Where("branch_id = ?", id).Update("branch_id", reassignToID).Error; err != nil {
			return err
		}

		// Спасаем дочерние ветки: поднимаем их на уровень выше к родителю удаляемой ветки
		if err := tx.Model(&model.Branch{}).Where("parent_id = ?", id).Update("parent_id", currentBranch.ParentID).Error; err != nil {
			return err
		}

		// Удаляем саму ветку
		return tx.Delete(&model.Branch{}, id).Error
	})
}
func (r *BranchRepository) GetByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error) {
	var branch model.Branch
	query := r.db.WithContext(ctx)

	// Защита от дурака: если передали какую-то дичь, сортируем по дате создания
	if sortBy != "created_at" && sortBy != "full_name" {
		sortBy = "created_at"
	}
	sortQuery := sortBy + " ASC" // ASC значит по возрастанию (от А до Я, от старых к новым)

	if includeEmployees {
		query = query.Preload("Flora", func(db *gorm.DB) *gorm.DB {
			return db.Order(sortQuery) // Вот тут используем нашу сортировку
		})
	}

	currentPath := "Children"
	for i := 1; i <= depth; i++ {
		query = query.Preload(currentPath)
		if includeEmployees {
			query = query.Preload(currentPath+".Flora", func(db *gorm.DB) *gorm.DB {
				return db.Order(sortQuery) // И тут тоже
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

// CheckCycle идет вверх по дереву от новой родительской ветки.
// Если по пути мы встретим саму себя (id), значит мы пытаемся засунуть ветку внутрь самой себя.
func (r *BranchRepository) CheckCycle(ctx context.Context, id uint, newParentID uint) (bool, error) {
	currentID := newParentID

	// Идем вверх, пока не упремся в корень (где родитель равен 0)
	for currentID != 0 {
		if currentID == id {
			return true, nil // Нашли цикл
		}

		var parentID *uint
		// Спрашиваем у базы только родителя текущей ветки
		err := r.db.WithContext(ctx).Model(&model.Branch{}).Select("parent_id").Where("id = ?", currentID).Scan(&parentID).Error
		if err != nil {
			return false, err
		}

		if parentID == nil {
			break // Дошли до самого верха, циклов нет
		}
		currentID = *parentID
	}

	return false, nil
}

// FindDuplicate проверяет, есть ли уже ветка с таким же именем в том же самом месте (у того же родителя).
func (r *BranchRepository) FindDuplicate(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	var count int64

	// Начинаем собирать запрос к базе: ищем ветку с таким же именем
	query := r.db.WithContext(ctx).Model(&model.Branch{}).Where("name = ?", name)

	// Проверяем родителя. Если parentID пустой (nil), значит ищем в самом корне сада
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	// Если мы обновляем ветку (excludeID не ноль), то просим базу не считать саму эту ветку дубликатом
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}

	// Считаем, сколько таких веток нашлось
	err := query.Count(&count).Error
	if err != nil {
		return false, err // Ой, ошибка базы данных
	}

	return count > 0, nil
}
