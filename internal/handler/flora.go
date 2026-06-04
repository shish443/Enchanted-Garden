// Enchanted-Garden/internal/handler/flora.go
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
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID ветви", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Сначала проверяем, существует ли вообще такая ветка в саду
	_, err = h.branchService.GetBranchByID(ctx, uint(id), 0, false, "")
	if err != nil {
		http.Error(w, "Ветка не найдена", http.StatusNotFound)
		return
	}

	var req model.PlantFloraReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" || utf8.RuneCountInString(req.FullName) > 200 {
		http.Error(w, "Имя должно быть от 1 до 200 символов", http.StatusBadRequest)
		return
	}
	req.Position = strings.TrimSpace(req.Position)
	if req.Position == "" || utf8.RuneCountInString(req.Position) > 200 {
		http.Error(w, "Позиция должна быть от 1 до 200 символов", http.StatusBadRequest)
		return
	}

	flora, err := h.service.CreateFlora(ctx, uint(id), &req)
	if err != nil {
		http.Error(w, "Ошибка при создании флоры: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(flora)
}
