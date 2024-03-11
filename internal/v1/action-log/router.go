package action_log

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListActionLog)
	r.Post("/", CreateActionLogEntry)

	r.Route("/{ActionLogID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Put("/", UpdateActionLog)
		r.Patch("/", PatchActionLog)
		r.Delete("/", DeleteActionLog)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ActionLogID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ActionLogID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		actionLog := &models.ActionLogEntry{ID: uint(ActionLogID)}
		if err = actionLog.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.AleKey{}, actionLog)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
