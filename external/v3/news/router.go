package news

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListNews)
	r.With(middleware.CanEditNews).Post("/", CreateNews)

	r.Route("/{NewsID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.CanEditNews).Put("/", UpdateNews)
		r.With(middleware.CanEditNews).Patch("/", PatchNews)
		r.With(middleware.CanEditNews).Delete("/", DeleteNews)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "NewsID")
		if id == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		NewsID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		news := &models.News{ID: uint(NewsID)}
		if err = news.Get(); err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.NewsKey{}, news)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
