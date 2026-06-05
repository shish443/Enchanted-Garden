// Enchanted-Garden/internal/service/flora.go
package service

import (
	"Enchanted-Garden/internal/model"
	"Enchanted-Garden/internal/repository"
	"context"
	"time"
)

type floraService struct {
	repo repository.FloraRepository
}

func NewFloraService(repo repository.FloraRepository) FloraService {
	return &floraService{repo: repo}
}

func (s *floraService) CreateFlora(ctx context.Context, branchID uint, req *model.PlantFloraReq) (*model.Flora, error) {
	var hiredAt *time.Time
	if req.HiredAt != nil {
		t := time.Time(*req.HiredAt)
		hiredAt = &t
	}

	f := &model.Flora{
		BranchID: branchID,
		FullName: req.FullName,
		Position: req.Position,
		HiredAt:  hiredAt,
	}

	if err := s.repo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}
