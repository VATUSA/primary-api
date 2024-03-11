package v1

import (
	"github.com/VATUSA/primary-api/internal/v1/document"
	facility_log "github.com/VATUSA/primary-api/internal/v1/facility-log"
	"github.com/VATUSA/primary-api/internal/v1/faq"
	"github.com/VATUSA/primary-api/internal/v1/feedback"
	"github.com/VATUSA/primary-api/internal/v1/news"
	"github.com/VATUSA/primary-api/internal/v1/notification"
	"github.com/VATUSA/primary-api/internal/v1/user"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router, cfg *config.Config) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/document", func(r chi.Router) {
			document.Router(r, cfg.S3)
		})

		r.Route("/faq", func(r chi.Router) {
			faq.Router(r)
		})

		r.Route("/facility-log", func(r chi.Router) {
			facility_log.Router(r)
		})

		r.Route("/feedback", func(r chi.Router) {
			feedback.Router(r)
		})

		r.Route("/news", func(r chi.Router) {
			news.Router(r)
		})

		r.Route("/notification", func(r chi.Router) {
			notification.Router(r)
		})

		r.Route("/user", func(r chi.Router) {
			user.Router(r)
		})
	})
}
