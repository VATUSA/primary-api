package v3

import (
	"github.com/VATUSA/primary-api/external/v3/event"
	"github.com/VATUSA/primary-api/external/v3/facility"
	"github.com/VATUSA/primary-api/external/v3/user"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router, cfg *config.Config) {
	r.Route("/v3", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			user.Router(r)
		})

		r.Route("/facility", func(r chi.Router) {
			facility.Router(r)
		})

		r.Get("/events", event.GetAllEvents)
	})
}
