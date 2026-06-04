// Enchanted-Garden/internal/handler/branch_test.go
package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBranchHandler_Create_BadJSON(t *testing.T) {
	handler := &BranchHandler{}

	// спецом делаем кривой текст чтобы сломать парсер
	badJSON := []byte(`{"name": "Новая ветка",,}`)

	req := httptest.NewRequest("POST", "/departments/", bytes.NewBuffer(badJSON))
	recorder := httptest.NewRecorder()

	handler.Create(recorder, req)

	// проверяем что статус равен 400 (ошибка клиента)
	// assert.Equal сам выведет красивую ошибку если что-то не совпадет
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
