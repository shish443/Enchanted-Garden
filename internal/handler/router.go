// Enchanted-Garden/internal/handler/router.go
package handler

import "net/http"

func SetupRouter(branchHandler *BranchHandler, floraHandler *FloraHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /departments/", branchHandler.Create)
	router.HandleFunc("GET /departments/{id}", branchHandler.GetByID)
	router.HandleFunc("PATCH /departments/{id}", branchHandler.Update)
	router.HandleFunc("DELETE /departments/{id}", branchHandler.Delete)

	router.HandleFunc("POST /departments/{id}/employees/", floraHandler.Create)

	return router
}
