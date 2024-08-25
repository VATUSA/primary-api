package facility

import (
	"context"
	facility_log "github.com/VATUSA/primary-api/external/v3/facility-log"
	"github.com/VATUSA/primary-api/external/v3/faq"
	"github.com/VATUSA/primary-api/external/v3/feedback"
	"github.com/VATUSA/primary-api/external/v3/news"
	"github.com/VATUSA/primary-api/external/v3/roster"
	roster_request "github.com/VATUSA/primary-api/external/v3/roster-request"
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

		r.With(middleware.NotGuest, middleware.CanEditFacility).Put("/", UpdateFacility)
		r.With(middleware.NotGuest, middleware.CanEditFacility).Patch("/", PatchFacility)
		r.With(middleware.NotGuest, middleware.CanEditFacility).Post("/reset-api-key", ResetApiKey)

		r.Route("/log", func(r chi.Router) {
			facility_log.Router(r)
		})

		r.Route("/faq", func(r chi.Router) {
			faq.Router(r)
		})

		r.Route("/feedback", func(r chi.Router) {
			feedback.Router(r)
		})

		r.Route("/news", func(r chi.Router) {
			news.Router(r)
		})

		r.Route("/roster", func(r chi.Router) {
			roster.Router(r)
		})

		r.Route("/roster-requests", func(r chi.Router) {
			roster_request.Router(r)
		})
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
