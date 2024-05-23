package facility

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(r chi.Router) {
	r.Get("/", GetFacilities)

	r.Route("/{FacilityID}", func(r chi.Router) {
		r.Use(Ctx)

		r.Get("/", GetFacility)
		r.With(middleware.CanEditFacility).Put("/", UpdateFacility)
		r.With(middleware.CanEditFacility).Patch("/", PatchFacility)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		facilityID := chi.URLParam(r, "FacilityID")
		if facilityID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		fac := &models.Facility{ID: constants.FacilityID(facilityID)}
		err := fac.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.FacilityKey{}, fac)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
