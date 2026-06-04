// Enchanted-Garden/logger/logger.go
package logger

import (
	"log/slog"
	"os"
)

// настраиваем вывод логов в консоль в красивом JSON
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, nil)

	// Делаем этот логгер главным для всего нашего проекта
	mainLogger := slog.New(handler)
	slog.SetDefault(mainLogger)
}
