package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
func New(log *slog.Logger, urlSaver URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("ropequest_id", middleware.GetReqID(r.Context())),
		)

		ailas := chi.URLParam(r, "alias")
		if ailas == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlSaver.GetURL(ailas)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "alias", ailas)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Info("failed to get url", "alias", ailas)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
