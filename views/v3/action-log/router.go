package action_log

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
	r.With(middleware.NotGuest, middleware.CanReadActionLog).Get("/", GetActionLog)
	r.With(middleware.NotGuest, middleware.CanEditActionLog).Post("/", CreateActionLogEntry)

	r.Route("/{ActionLogID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.NotGuest, middleware.CanEditActionLog).Patch("/", PatchActionLogEntry)
		r.With(middleware.NotGuest, middleware.CanEditActionLog).Put("/", UpdateActionLogEntry)
		r.With(middleware.NotGuest, middleware.CanEditActionLog).Delete("/", DeleteActionLogEntry)
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
