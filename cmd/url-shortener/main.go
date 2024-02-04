package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/urldelete"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
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
	// добавляет идентификатор каждому запросу
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP) // ip пользователя

	// лог запросов из коробки. Проблема, что у нас свой логгер
	//router.Use(middleware.Logger)
	// своя реализация логгера для middleware
	router.Use(mwLogger.New(log))

	router.Use(middleware.Recoverer) // приложение не падает при плохом запросе
	router.Use(middleware.URLFormat) // можно писать в хендлере красивые урлы типа /articles/{id}. И обращаться по {id}

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cnf.HTTPServer.User: cnf.HTTPServer.Password,
			//cnf.HTTPServer.User: cnf.HTTPServer.Password, // может добавить других пользователей
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", urldelete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cnf.Address))

	// TODO run server
	srv := http.Server{
		Addr:         cnf.Address,
		Handler:      router,
		ReadTimeout:  cnf.Timeout,
		WriteTimeout: cnf.Timeout,
		IdleTimeout:  cnf.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

// лог (вид и уровень) зависит от окружения: dev, prod и тд
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
		//log = slog.New(
		//    slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		//)
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
