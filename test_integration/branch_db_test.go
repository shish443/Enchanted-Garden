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
		t.Skip("database connection unavailable, skipping test: ", err)
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

	testBranch := &model.Branch{
		Name: "Тестовая Аллея",
	}
	require.NoError(t, repo.Create(ctx, testBranch))
	assert.NotEqual(t, 0, testBranch.ID)

	found, err := repo.GetByID(ctx, testBranch.ID, 1, false, "")
	require.NoError(t, err)
	assert.Equal(t, "Тестовая Аллея", found.Name)

	found.Name = "Магическая Аллея"
	require.NoError(t, repo.Update(ctx, found))

	updated, err := repo.GetByID(ctx, testBranch.ID, 1, false, "")
	require.NoError(t, err)
	assert.Equal(t, "Магическая Аллея", updated.Name)

	require.NoError(t, repo.DeleteCascade(ctx, testBranch.ID))

	_, err = repo.GetByID(ctx, testBranch.ID, 1, false, "")
	require.Error(t, err)
}
