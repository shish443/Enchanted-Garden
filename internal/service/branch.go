// Enchanted-Garden/internal/service/branch.go

package service

import (
	"Enchanted-Garden/internal/model"
	"Enchanted-Garden/internal/repository"
	"context"
	"errors"
	"fmt"
)

type branchService struct {
	repo repository.BranchRepository
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchService{repo: repo}
}

func (s *branchService) CreateBranch(ctx context.Context, req *model.CreateBranchReq) (*model.Branch, error) {
	// Проверяем уникальность имени у этого родителя
	isDuplicate, err := s.repo.FindDuplicate(ctx, req.ParentID, req.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("check duplicate failed: %w", err)
	}
	if isDuplicate {
		return nil, errors.New("branch name must be unique within the same parent")
	}

	b := &model.Branch{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	err = s.repo.Create(ctx, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *branchService) GetBranchByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}
	return s.repo.GetByID(ctx, id, depth, includeEmployees, sortBy)
}

func (s *branchService) UpdateBranch(ctx context.Context, id uint, req *model.UpdateBranchReq) (*model.Branch, error) {
	b, err := s.repo.GetByID(ctx, id, 0, false, "")
	if err != nil {
		return nil, err
	}

	finalName := b.Name
	if req.Name != nil {
		finalName = *req.Name
	}

	finalParentID := b.ParentID
	if req.ParentID != nil {
		// Если нам передали поле, берем то, что внутри (может быть nil, а может быть число)
		finalParentID = *req.ParentID
	}

	// Проверяем уникальность имени при обновлении
	isDuplicate, _ := s.repo.FindDuplicate(ctx, finalParentID, finalName, id)
	if isDuplicate {
		return nil, errors.New("branch name must be unique within the same parent")
	}

	if req.Name != nil {
		b.Name = *req.Name
	}

	// Если нам прислали запрос на изменение родителя
	if req.ParentID != nil {

		// Если прислали конкретную другую ветку (не null)
		if *req.ParentID != nil {
			newParentID := **req.ParentID // достаем число

			if newParentID == id {
				return nil, errors.New("нельзя сделать ветку родителем самой себя")
			}

			// Проверяем, не зациклили ли мы наш сад
			hasCycle, err := s.repo.CheckCycle(ctx, id, newParentID)
			if err != nil {
				return nil, err
			}
			if hasCycle {
				return nil, errors.New("cycle detected")
			}
		}

		// Сохраняем нового родителя (если прислали null, тут запишется nil, и ветка станет корневой)
		b.ParentID = *req.ParentID
	}

	err = s.repo.Update(ctx, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *branchService) DeleteBranch(ctx context.Context, id uint, mode string, reassignToID uint) error {
	if mode == "cascade" {
		return s.repo.DeleteCascade(ctx, id)
	}

	if mode == "reassign" {
		// Сначала пробуем найти ветку, куда собираемся отправить сотрудников
		_, err := s.repo.GetByID(ctx, reassignToID, 1, false, "")
		if err != nil {
			// Ветки нет, отдаем ошибку
			return errors.New("target department not found")
		}

		return s.repo.DeleteReassign(ctx, id, reassignToID)
	}

	return errors.New("неизвестный режим удаления")
}
