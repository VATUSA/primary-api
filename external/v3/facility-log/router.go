package facility_log

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
	r.With(middleware.CanViewFacilityLog).Get("/", ListFacilityLog)

	r.Route("/{FacilityLogID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.CanEditFacilityLog).Put("/", UpdateFacilityLog)
		r.With(middleware.CanEditFacilityLog).Patch("/", PatchFacilityLog)
		r.With(middleware.CanEditFacilityLog).Delete("/", DeleteFacilityLog)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		facilityLogID := chi.URLParam(r, "FacilityLogID")
		if facilityLogID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		facilityLogIDInt, err := strconv.Atoi(facilityLogID)
		if err != nil {
			utils.Render(w, r, utils.ErrInvalidRequest(err))
			return
		}

		fac := &models.FacilityLogEntry{ID: uint(facilityLogIDInt)}
		err = fac.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.FacilityLogKey{}, fac)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
