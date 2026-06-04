// Enchanted-Garden/internal/handler/router.go
package handler

import "net/http"

// SetupRouter настраивает все пути для сада.Принимает готовые обработчики и возвращает роутер
func SetupRouter(branchHandler *BranchHandler, floraHandler *FloraHandler) *http.ServeMux {
	// Создаем новый роутер
	router := http.NewServeMux()

	// Пути для веток (подразделений)
	router.HandleFunc("POST /departments/", branchHandler.Create)
	router.HandleFunc("GET /departments/{id}", branchHandler.GetByID)
	router.HandleFunc("PATCH /departments/{id}", branchHandler.Update)
	router.HandleFunc("DELETE /departments/{id}", branchHandler.Delete)

	// Пути для флоры (сотрудников)
	router.HandleFunc("POST /departments/{id}/employees/", floraHandler.Create)

	return router
}
