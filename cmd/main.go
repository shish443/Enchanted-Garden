//Enchanted Garden/cmd/main.go

package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	//отображение джейсона к консоль в норм виде
	logger := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(logger)
	slog.Info("start garden api")
	//фактически наш дворечкий, решает что куда отправить
	router := http.NewServeMux()
	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	//start setver
	slog.Info("Server run in 8080 port")
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("Server dont start", "error", err)
		os.Exit(1)
	}
}
