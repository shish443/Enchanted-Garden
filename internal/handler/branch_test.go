// Enchanted-Garden/internal/handler/branch_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Enchanted-Garden/internal/model"

	"github.com/stretchr/testify/assert"
)

type stubBranchService struct {
	onCreate func(ctx context.Context, req *model.CreateBranchReq) (*model.Branch, error)
}

func (s *stubBranchService) CreateBranch(ctx context.Context, req *model.CreateBranchReq) (*model.Branch, error) {
	if s.onCreate != nil {
		return s.onCreate(ctx, req)
	}
	return &model.Branch{}, nil
}

func (s *stubBranchService) GetBranchByID(ctx context.Context, id uint, depth int, inc bool, sort string) (*model.Branch, error) {
	return nil, nil
}

func (s *stubBranchService) UpdateBranch(ctx context.Context, id uint, req *model.UpdateBranchReq) (*model.Branch, error) {
	return nil, nil
}

func (s *stubBranchService) DeleteBranch(ctx context.Context, id uint, mode string, reassignToID uint) error {
	return nil
}

func TestBranchHandler_Create_Success(t *testing.T) {
	fakeService := &stubBranchService{
		onCreate: func(ctx context.Context, req *model.CreateBranchReq) (*model.Branch, error) {
			return &model.Branch{ID: 10, Name: req.Name}, nil
		},
	}
	h := NewBranchHandler(fakeService)

	reqBody := model.CreateBranchReq{Name: "Волшебная Аллея"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/departments/", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	h.Create(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var res model.Branch
	err := json.NewDecoder(recorder.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, uint(10), res.ID)
	assert.Equal(t, "Волшебная Аллея", res.Name)
}

func TestBranchHandler_Create_ValidationError(t *testing.T) {
	h := NewBranchHandler(&stubBranchService{})

	reqBody := model.CreateBranchReq{Name: "   "}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/departments/", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	h.Create(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "name cannot be empty")
}
