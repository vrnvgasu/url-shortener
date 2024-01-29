package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lig/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	endDev   = "dev"
	endProd  = "prod"
)

func main() {
	// TODO init config: cleanenv - может читать из разных источников
	cnf := config.MustLoad()
	fmt.Println(cnf)

	// TODO init logger: slog - в ядре с версии 1.21
	log := setupLogger(cnf.Env)
	//log = log.With("env", cnf.Env) // добавляем параметр env ко всем логам
	log.Info("starting application", slog.String("env", cnf.Env))
	log.Debug("debug messages are enabled")

	// TODO init storage: sqlite
	storage, err := sqlite.New(cnf.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO init router: chi, "chi render"
	router := chi.NewRouter()
	_ = router

	// TODO run server
}

// лог (вид и уровень) зависит от окружения: dev, prod и тд
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case endDev: // для dev стенда
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case endProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
