package disciplinary_log

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
	r.With(middleware.NotGuest, middleware.CanReadDisciplinaryLog).Get("/", GetDisciplinaryLog)
	r.With(middleware.NotGuest, middleware.CanEditDisciplinaryLog).Post("/", CreateDisciplinaryLogEntry)

	r.Route("/{DisciplinaryLogID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.NotGuest, middleware.CanEditDisciplinaryLog).Patch("/", PatchDisciplinaryLogEntry)
		r.With(middleware.NotGuest, middleware.CanEditDisciplinaryLog).Put("/", UpdateDisciplinaryLogEntry)
		r.With(middleware.NotGuest, middleware.CanEditDisciplinaryLog).Delete("/", DeleteDisciplinaryLogEntry)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "DisciplinaryLogID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		DisciplinaryLogID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		disciplinaryLog := &models.DisciplinaryLogEntry{ID: uint(DisciplinaryLogID)}
		if err = disciplinaryLog.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.DleKey{}, disciplinaryLog)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
