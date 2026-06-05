// Enchanted-Garden/internal/service/branch.go
package service

import (
	"Enchanted-Garden/internal/model"
	"Enchanted-Garden/internal/repository"
	"context"
	"errors"
)

var (
	ErrBranchDuplicate = errors.New("branch name must be unique within the same parent")
	ErrBranchCycle     = errors.New("cycle detected")
	ErrBranchNotFound  = errors.New("branch not found")
)

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) CreateBranch(ctx context.Context, req *model.CreateBranchReq) (*model.Branch, error) {
	isDup, err := s.repo.FindDuplicate(ctx, req.ParentID, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if isDup {
		return nil, ErrBranchDuplicate
	}

	branch := &model.Branch{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if err := s.repo.Create(ctx, branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *branchService) UpdateBranch(ctx context.Context, id uint, req *model.UpdateBranchReq) (*model.Branch, error) {
	branch, err := s.repo.GetByID(ctx, id, 0, false, "")
	if err != nil {
		return nil, ErrBranchNotFound
	}

	if req.Name != nil || req.ParentID != nil {
		name := branch.Name
		if req.Name != nil {
			name = *req.Name
		}

		parentID := branch.ParentID
		if req.ParentID != nil {
			if *req.ParentID == 0 {
				parentID = nil
			} else {
				parentID = req.ParentID
			}
		}

		isDup, err := s.repo.FindDuplicate(ctx, parentID, name, id)
		if err != nil {
			return nil, err
		}
		if isDup {
			return nil, ErrBranchDuplicate
		}
	}

	if req.ParentID != nil {
		if *req.ParentID != 0 {
			if *req.ParentID == id {
				return nil, errors.New("cannot make a branch a parent of itself")
			}

			hasCycle, err := s.repo.CheckCycle(ctx, id, *req.ParentID)
			if err != nil {
				return nil, err
			}
			if hasCycle {
				return nil, ErrBranchCycle
			}
			branch.ParentID = req.ParentID
		} else {
			branch.ParentID = nil
		}
	}

	if req.Name != nil {
		branch.Name = *req.Name
	}

	if err := s.repo.Update(ctx, branch); err != nil {
		return nil, err
	}

	return branch, nil
}

func (s *branchService) DeleteBranch(ctx context.Context, id uint, mode string, reassignToID uint) error {
	_, err := s.repo.GetByID(ctx, id, 0, false, "")
	if err != nil {
		return ErrBranchNotFound
	}

	if mode == "reassign" {
		if reassignToID == 0 {
			return errors.New("reassign id is required")
		}
		return s.repo.DeleteReassign(ctx, id, reassignToID)
	}

	return s.repo.DeleteCascade(ctx, id)
}

func (s *branchService) GetBranchByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}
	branch, err := s.repo.GetByID(ctx, id, depth, includeEmployees, sortBy)
	if err != nil {
		return nil, ErrBranchNotFound
	}
	return branch, nil
}
