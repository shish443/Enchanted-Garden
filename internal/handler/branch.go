package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"Enchanted-Garden/internal/model"
	"Enchanted-Garden/internal/service"
)

type BranchHandler struct {
	service service.BranchService
}

func NewBranchHandler(s service.BranchService) *BranchHandler {
	return &BranchHandler{service: s}
}

func (h *BranchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBranchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || utf8.RuneCountInString(req.Name) > 200 {
		http.Error(w, "name cannot be empty or exceed 200 characters", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	branch, err := h.service.CreateBranch(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "branch name already exists under this parent", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create branch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(branch)
}

func (h *BranchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid branch ID", http.StatusBadRequest)
		return
	}

	depth := 1
	if depthStr := r.URL.Query().Get("depth"); depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil {
			depth = d
		}
	}
	if depth < 1 {
		depth = 1
	} else if depth > 5 {
		depth = 5
	}

	includeEmployees := true
	if r.URL.Query().Get("include_employees") == "false" {
		includeEmployees = false
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	branch, err := h.service.GetBranchByID(r.Context(), uint(id), depth, includeEmployees, sortBy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "branch not found"}`, http.StatusNotFound)
			return
		}

		slog.Error("failed to get branch", "error", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func (h *BranchHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid branch ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateBranchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" || utf8.RuneCountInString(trimmed) > 200 {
			http.Error(w, "branch name cannot be empty or exceed 200 characters", http.StatusBadRequest)
			return
		}
		req.Name = &trimmed
	}

	if req.ParentID != nil {
		if *req.ParentID == uint(id) {
			http.Error(w, "branch cannot be a parent of itself", http.StatusBadRequest)
			return
		}
	}

	branch, err := h.service.UpdateBranch(r.Context(), uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "branch not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "branch name already exists under this parent", http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "cycle") {
			http.Error(w, "circular dependency detected: cannot move branch under its own subtree", http.StatusConflict)
			return
		}
		http.Error(w, "failed to update branch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func (h *BranchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid branch ID", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode != "cascade" && mode != "reassign" {
		http.Error(w, "mode parameter must be 'cascade' or 'reassign'", http.StatusBadRequest)
		return
	}

	var reassignToID uint
	if mode == "reassign" {
		reassignStr := r.URL.Query().Get("reassign_to_department_id")
		if reassignStr == "" {
			http.Error(w, "reassign_to_department_id is required when mode is 'reassign'", http.StatusBadRequest)
			return
		}
		reassignTo, err := strconv.ParseUint(reassignStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid reassign branch ID", http.StatusBadRequest)
			return
		}
		reassignToID = uint(reassignTo)

		if reassignToID == uint(id) {
			http.Error(w, "cannot reassign entities to the branch being deleted", http.StatusBadRequest)
			return
		}

		if _, err = h.service.GetBranchByID(r.Context(), reassignToID, 1, false, ""); err != nil {
			http.Error(w, "target reassignment branch not found", http.StatusNotFound)
			return
		}
	}

	if err = h.service.DeleteBranch(r.Context(), uint(id), mode, reassignToID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "branch not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete branch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
