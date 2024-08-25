package roster

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
	r.With(middleware.NotGuest, middleware.CanEditRoster).Post("/", CreateRoster)
	r.Get("/", GetRosterByFacility)
	r.Route("/{RosterID}", func(r chi.Router) {
		r.Use(Ctx)
		r.With(middleware.NotGuest, middleware.CanEditRoster).Delete("/", DeleteRoster)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "RosterID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		RosterID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		roster := &models.Roster{ID: uint(RosterID)}
		if err = roster.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.RosterKey{}, roster)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
