package feedback

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListFeedback)
	r.Post("/", CreateFeedback)

	r.Route("/{FeedbackID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Get("/", GetFeedback)
		r.Put("/", UpdateFeedback)
		r.Patch("/", PatchFeedback)
		r.Delete("/", DeleteFeedback)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "FeedbackID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		FeedbackID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		feedback := &models.Feedback{ID: uint(FeedbackID)}
		if err = feedback.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.FeedbackKey{}, feedback)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
