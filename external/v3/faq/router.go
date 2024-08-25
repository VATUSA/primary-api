package faq

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
	r.Get("/", ListFAQ)

	r.With(middleware.NotGuest, middleware.CanEditFAQ).Post("/", CreateFAQ)

	r.Route("/{FAQID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.NotGuest, middleware.CanEditFAQ).Put("/", UpdateFAQ)
		r.With(middleware.NotGuest, middleware.CanEditFAQ).Patch("/", PatchFAQ)
		r.With(middleware.NotGuest, middleware.CanEditFAQ).Delete("/", DeleteFAQ)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "FAQID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		faqID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		faq := &models.FAQ{ID: uint(faqID)}
		if err := faq.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.FAQKey{}, faq)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
