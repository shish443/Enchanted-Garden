package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"Enchanted-Garden/internal/model"
	"Enchanted-Garden/internal/service"
)

type FloraHandler struct {
	service       service.FloraService
	branchService service.BranchService
}

func NewFloraHandler(s service.FloraService, b service.BranchService) *FloraHandler {
	return &FloraHandler{service: s, branchService: b}
}

func (h *FloraHandler) Create(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid branch ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if _, err = h.branchService.GetBranchByID(ctx, uint(id), 0, false, ""); err != nil {
		http.Error(w, "target branch not found", http.StatusNotFound)
		return
	}

	var req model.PlantFloraReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" || utf8.RuneCountInString(req.FullName) > 200 {
		http.Error(w, "full name must be between 1 and 200 characters", http.StatusBadRequest)
		return
	}

	req.Position = strings.TrimSpace(req.Position)
	if req.Position == "" || utf8.RuneCountInString(req.Position) > 200 {
		http.Error(w, "position must be between 1 and 200 characters", http.StatusBadRequest)
		return
	}

	flora, err := h.service.CreateFlora(ctx, uint(id), &req)
	if err != nil {
		http.Error(w, "failed to create flora: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flora)
}
