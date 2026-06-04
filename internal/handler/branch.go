// Enchanted-Garden/internal/handler/branch.go
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
		http.Error(w, "неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	// Считаем именно символы, а не байты, чтобы кириллица не ломала длину
	if req.Name == "" || utf8.RuneCountInString(req.Name) > 200 {
		http.Error(w, "Имя не должно быть пустым или длиннее 200 символов", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	branch, err := h.service.CreateBranch(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "Ветка с таким именем уже существует у этого родителя", http.StatusBadRequest)
			return
		}
		http.Error(w, "Ошибка создания ветви: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(branch)
}

func (h *BranchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID ветви", http.StatusBadRequest)
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

	// Читаем из ссылки, нужно ли нам доставать сотрудников
	queryParams := r.URL.Query()
	includeEmployeesStr := queryParams.Get("include_employees")

	includeEmployees := true
	if includeEmployeesStr == "false" {
		includeEmployees = false
	}

	sortBy := r.URL.Query().Get("sort_by") // Читаем параметр из ссылки
	if sortBy == "" {
		sortBy = "created_at" // Ставим по умолчанию
	}

	branch, err := h.service.GetBranchByID(r.Context(), uint(id), depth, includeEmployees, sortBy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Ветвь не найдена"}`, http.StatusNotFound)
			return
		}

		// Пишем реальную ошибку в консоль для себя
		slog.Error("Ошибка получения ветви", "error", err)
		// Клиенту отдаем безопасный текст без деталей базы данных
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func (h *BranchHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID ветви", http.StatusBadRequest)
		return
	}

	var req model.UpdateBranchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" || utf8.RuneCountInString(trimmed) > 200 {
			http.Error(w, "Имя ветви не должно быть пустым или длиннее 200 символов", http.StatusBadRequest)
			return
		}
		req.Name = &trimmed
	}

	// Если parent_id вообще прислали в JSON
	if req.ParentID != nil {
		// Если прислали конкретное число, а не null
		if *req.ParentID != nil {
			// Достаем само число из двойного указателя
			if **req.ParentID == uint(id) {
				http.Error(w, "Ветвь не может быть родителем самой себя", http.StatusBadRequest)
				return
			}
		}
	}

	branch, err := h.service.UpdateBranch(r.Context(), uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Ветвь не найдена", http.StatusNotFound)
			return
		}

		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "Ветка с таким именем уже существует у этого родителя", http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "cycle") {
			http.Error(w, "Нельзя переместить ветвь внутрь своего поддерева (цикл)", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка обновления ветви: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func (h *BranchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID ветви", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode != "cascade" && mode != "reassign" {
		http.Error(w, "параметр mode должен быть 'cascade' или 'reassign'", http.StatusBadRequest)
		return
	}

	var reassignToID uint
	if mode == "reassign" {
		reassignStr := r.URL.Query().Get("reassign_to_department_id")
		if reassignStr == "" {
			http.Error(w, "Для режима 'reassign' параметр reassign_to_department_id обязателен", http.StatusBadRequest)
			return
		}
		reassignTo, err := strconv.ParseUint(reassignStr, 10, 64)
		if err != nil {
			http.Error(w, "Неверный ID ветви для переназначения флоры", http.StatusBadRequest)
			return
		}
		reassignToID = uint(reassignTo)

		if reassignToID == uint(id) {
			http.Error(w, "Нельзя переназначить флору в удаляемую ветвь", http.StatusBadRequest)
			return
		}

		_, err = h.service.GetBranchByID(r.Context(), reassignToID, 1, false, "")
		if err != nil {
			// Если произошла ошибка (ветки нет), сразу ругаемся и останавливаем работу
			http.Error(w, "Ветка для пересадки флоры не найдена в саду", http.StatusNotFound)
			return
		}
	}

	err = h.service.DeleteBranch(r.Context(), uint(id), mode, reassignToID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Ветвь не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при удалении ветви: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
