// Enchanted-Garden/internal/service/branch_test.go
package service

import (
	"context"
	"testing"

	"Enchanted-Garden/internal/model"

	"github.com/stretchr/testify/assert"
)

type dummyBranchRepo struct {
	lastDepth int
}

func (d *dummyBranchRepo) Create(ctx context.Context, branch *model.Branch) error { return nil }

func (d *dummyBranchRepo) GetByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error) {
	d.lastDepth = depth
	return &model.Branch{ID: id, Name: "Test Branch"}, nil
}

func (d *dummyBranchRepo) Update(ctx context.Context, branch *model.Branch) error { return nil }
func (d *dummyBranchRepo) DeleteCascade(ctx context.Context, id uint) error       { return nil }
func (d *dummyBranchRepo) DeleteReassign(ctx context.Context, id uint, reassignToID uint) error {
	return nil
}

func (d *dummyBranchRepo) FindDuplicate(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	return false, nil
}

func (d *dummyBranchRepo) CheckCycle(ctx context.Context, id uint, newParentID uint) (bool, error) {
	return false, nil
}

type spyBranchRepo struct {
	dummyBranchRepo
	mockFindDuplicate func(parentID *uint, name string) (bool, error)
}

func (s *spyBranchRepo) FindDuplicate(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	if s.mockFindDuplicate != nil {
		return s.mockFindDuplicate(parentID, name)
	}
	return false, nil
}

func TestGetBranchByID_MaxDepthLimit(t *testing.T) {
	repo := &dummyBranchRepo{}
	s := NewBranchService(repo)

	_, err := s.GetBranchByID(context.Background(), 1, 10, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastDepth != 5 {
		t.Errorf("expected depth to be capped at 5, got %d", repo.lastDepth)
	}
}

func TestCreateBranch_DuplicateNameError(t *testing.T) {
	repo := &spyBranchRepo{
		mockFindDuplicate: func(parentID *uint, name string) (bool, error) {
			return true, nil
		},
	}
	s := NewBranchService(repo)

	req := &model.CreateBranchReq{Name: "Дубликат"}
	_, err := s.CreateBranch(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "branch name must be unique")
}
