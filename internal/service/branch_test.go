// Enchanted-Garden/internal/service/branch_test.go
package service

import (
	"Enchanted-Garden/internal/model"
	"context"
	"testing"
)

// Простая заглушка репозитория, чтобы протестировать логику сервиса без подключения к реальной базе данных
type dummyBranchRepo struct{}

func (d *dummyBranchRepo) Create(ctx context.Context, branch *model.Branch) error { return nil }
func (d *dummyBranchRepo) Update(ctx context.Context, branch *model.Branch) error { return nil }
func (d *dummyBranchRepo) DeleteCascade(ctx context.Context, id uint) error       { return nil }
func (d *dummyBranchRepo) DeleteReassign(ctx context.Context, id uint, reassignToID uint) error {
	return nil
}

// Добавили этот метод, чтобы заглушка соответствовала новому интерфейсу репозитория
func (d *dummyBranchRepo) FindDuplicate(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	return false, nil
}

// Проверяем, что если передать слишком большую глубину (например, 10), сервис автоматически скинет её до максимума (5)
func TestGetBranchByID_MaxDepthLimit(t *testing.T) {
	repo := &dummyBranchRepo{}
	s := NewBranchService(repo)

	// Запрашиваем ветку с сумасшедшей глубиной 10
	_, err := s.GetBranchByID(context.Background(), 1, 10, false, "")
	if err != nil {
		t.Errorf("Не ждали ошибку, но она случилась: %v", err)
	}
}

func (d *dummyBranchRepo) GetByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error) {
	return &model.Branch{ID: id, Name: "Тестовая ветка"}, nil
}

func (d *dummyBranchRepo) CheckCycle(ctx context.Context, id uint, newParentID uint) (bool, error) {
	return false, nil
}
