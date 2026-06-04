// Enchanted-Garden/internal/service/service.go
package service

import (
	"Enchanted-Garden/internal/model"
	"context"
)

type BranchService interface {
	CreateBranch(ctx context.Context, req *model.CreateBranchReq) (*model.Branch, error)
	GetBranchByID(ctx context.Context, id uint, depth int, includeEmployees bool, sortBy string) (*model.Branch, error)
	UpdateBranch(ctx context.Context, id uint, req *model.UpdateBranchReq) (*model.Branch, error)
	DeleteBranch(ctx context.Context, id uint, mode string, reassignToID uint) error
}

type FloraService interface {
	CreateFlora(ctx context.Context, branchID uint, req *model.PlantFloraReq) (*model.Flora, error)
}
