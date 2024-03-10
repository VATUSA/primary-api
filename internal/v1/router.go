package v1

import (
	action_log "github.com/VATUSA/primary-api/internal/v1/action-log"
	disciplinary_log "github.com/VATUSA/primary-api/internal/v1/disciplinary-log"
	"github.com/VATUSA/primary-api/internal/v1/document"
	facility_log "github.com/VATUSA/primary-api/internal/v1/facility-log"
	"github.com/VATUSA/primary-api/internal/v1/faq"
	"github.com/VATUSA/primary-api/internal/v1/feedback"
	"github.com/VATUSA/primary-api/internal/v1/news"
	"github.com/VATUSA/primary-api/internal/v1/notification"
	rating_change "github.com/VATUSA/primary-api/internal/v1/rating-change"
	roster_request "github.com/VATUSA/primary-api/internal/v1/roster-request"
	"github.com/VATUSA/primary-api/internal/v1/user"
	user_role "github.com/VATUSA/primary-api/internal/v1/user-role"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router, cfg *config.Config) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/action-log", func(r chi.Router) {
			action_log.Router(r)
		})

		r.Route("/disciplinary-log", func(r chi.Router) {
			disciplinary_log.Router(r)
		})

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

		r.Route("/rating-change", func(r chi.Router) {
			rating_change.Router(r)
		})

		r.Route("/roster-request", func(r chi.Router) {
			roster_request.Router(r)
		})

		r.Route("/user", func(r chi.Router) {
			user.Router(r)
		})

		r.Route("/user-role", func(r chi.Router) {
			user_role.Router(r)
		})
	})
}
