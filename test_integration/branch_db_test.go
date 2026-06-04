// Enchanted-Garden/test_integration/branch_db_test.go
package testintegration

import (
	"context"
	"os"
	"testing"

	"Enchanted-Garden/internal/model"
	postgresRepo "Enchanted-Garden/internal/repository/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/garden_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("Пропускаем тест, база сейчас недоступна: ", err)
		return nil
	}

	db.AutoMigrate(&model.Branch{}, &model.Flora{})
	return db
}

func TestBranchRepository_FullLifeCycle(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := postgresRepo.NewBranchRepository(db)
	ctx := context.Background()

	// 1. Создание
	testBranch := &model.Branch{
		Name: "Тестовая Аллея",
	}
	err := repo.Create(ctx, testBranch)
	// require.NoError сразу останавливает тест если есть ошибка
	require.NoError(t, err)
	// проверяем что база выдала нам нормальный айдишник (не ноль)
	assert.NotEqual(t, 0, testBranch.ID)

	// 2. Чтение
	found, err := repo.GetByID(ctx, testBranch.ID, 1, false, "")
	require.NoError(t, err)
	assert.Equal(t, "Тестовая Аллея", found.Name)

	// 3. Обновление
	found.Name = "Магическая Аллея"
	err = repo.Update(ctx, found)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, testBranch.ID, 1, false, "")
	require.NoError(t, err)
	assert.Equal(t, "Магическая Аллея", updated.Name)

	// 4. Удаление
	err = repo.DeleteCascade(ctx, testBranch.ID)
	require.NoError(t, err)

	// 5. Проверяем что ветки точно больше нет
	_, err = repo.GetByID(ctx, testBranch.ID, 1, false, "")
	// тут мы ПРЯМ ЖДЕМ ошибку, так как удалили ветку
	require.Error(t, err)
}
