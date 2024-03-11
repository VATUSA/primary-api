package rating_change

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListRatingChanges)
	r.Post("/", CreateRatingChange)

	r.Route("/{RatingChangeID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Put("/", UpdateRatingChange)
		r.Patch("/", PatchRatingChange)
		r.Delete("/", DeleteRatingChange)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "RatingChangeID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		RatingChangeID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ratingChange := &models.RatingChange{ID: uint(RatingChangeID)}
		if err = ratingChange.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "ratingChange", ratingChange)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
