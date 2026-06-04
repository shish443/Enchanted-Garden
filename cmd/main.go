// Enchanted-Garden/cmd/main.go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Enchanted-Garden/internal/config"
	"Enchanted-Garden/internal/handler"
	"Enchanted-Garden/internal/repository/postgres"
	"Enchanted-Garden/internal/service"
	"Enchanted-Garden/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger.Init()
	slog.Info("Запускаем наш сад")

	cfg := config.Load()

	db, err := gorm.Open(gormpostgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		slog.Error("Не смогли подключиться к базе", "error", err)
		os.Exit(1)
	}
	slog.Info("База данных подключена")

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Не смогли получить sql.DB из GORM", "error", err)
		os.Exit(1)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		slog.Error("Не смогли настроить goose на работу с postgres", "error", err)
		os.Exit(1)
	}

	err = goose.Up(sqlDB, "db/migrations")
	if err != nil {
		slog.Error("Ошибка во время выполнения миграций", "error", err)
		os.Exit(1)
	}
	slog.Info("Миграции успешно применены!")

	branchRepo := postgres.NewBranchRepository(db)
	floraRepo := postgres.NewFloraRepository(db)

	branchService := service.NewBranchService(branchRepo)
	floraService := service.NewFloraService(floraRepo)

	branchHandler := handler.NewBranchHandler(branchService)
	floraHandler := handler.NewFloraHandler(floraService, branchService)

	router := handler.SetupRouter(branchHandler, floraHandler)

	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// создаем сервер руками, чтобы им можно было управлять
	srv := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}

	// запускаем сервер в фоне, чтобы код пошел дальше
	go func() {
		slog.Info("Сервер слушает порт " + cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Сервер упал", "error", err)
		}
	}()

	// ждем когда нажмут ctrl+c
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // тут программа зависает и ждет сигнал

	slog.Info("Выключаем сервер...")

	// даем 5 секунд на завершение того что уже начало работать
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Ошибка выключения", "error", err)
	}
	slog.Info("Закрываем соединения с базой данных...")
	if err := sqlDB.Close(); err != nil {
		slog.Error("Ошибка при закрытии БД", "error", err)
	}
	slog.Info("Всё остановлено")
}
