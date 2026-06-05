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
	slog.Info("starting application")

	cfg := config.Load()

	db, err := gorm.Open(gormpostgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB", "error", err)
		os.Exit(1)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("failed to set goose dialect", "error", err)
		os.Exit(1)
	}

	if err := goose.Up(sqlDB, "db/migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("database migrations applied successfully")

	branchRepo := postgres.NewBranchRepository(db)
	floraRepo := postgres.NewFloraRepository(db)

	branchService := service.NewBranchService(branchRepo)
	floraService := service.NewFloraService(floraRepo)

	branchHandler := handler.NewBranchHandler(branchService)
	floraHandler := handler.NewFloraHandler(floraService, branchService)

	router := handler.SetupRouter(branchHandler, floraHandler)
	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	srv := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}

	go func() {
		slog.Info("http server started", "port", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
	}

	slog.Info("application stopped completely")
}
