package facility_log

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListFacilityLog)
	r.Post("/", CreateFacilityLogEntry)

	r.Route("/{FacilityLogID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Get("/", GetFacilityLog)
		r.Put("/", UpdateFacilityLog)
		r.Patch("/", PatchFacilityLog)
		r.Delete("/", DeleteFacilityLog)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "FacilityLogID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		FacilityLogID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		facilityLog := &models.FacilityLogEntry{ID: uint(FacilityLogID)}
		if err = facilityLog.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.FacilityLogKey{}, facilityLog)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
