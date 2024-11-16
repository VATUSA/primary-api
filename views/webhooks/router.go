package v3

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/views/webhooks/vatsim"
	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router, cfg *config.Config) {
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/vatsim", vatsim.ProcessMemberWebhook)
	})
}
