package facility

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(r chi.Router) {
	r.Get("/", GetFacilities)

	r.Route("/{FacilityID}", func(r chi.Router) {
		r.Use(Ctx)

		r.Get("/", GetFacility)
		r.Put("/", UpdateFacility)
		r.Patch("/", PatchFacility)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		facilityID := chi.URLParam(r, "FacilityID")
		if facilityID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		fac := &models.Facility{ID: facilityID}
		err := fac.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.FacilityLogKey{}, fac)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
