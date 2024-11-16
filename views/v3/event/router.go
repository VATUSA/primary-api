package event

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func EventRouter(r chi.Router) {
	r.Get("/", GetEvents)
	r.Get("/previous", GetEventsPrevious)
	r.With(middleware.NotGuest, middleware.CanEditEvent).Post("/", CreateEvent)

	r.Route("/{EventID}", func(r chi.Router) {
		r.Use(eventCtx)

		r.Get("/", GetEvent)
		r.With(middleware.NotGuest, middleware.CanEditEvent).Put("/", UpdateEvent)
		r.With(middleware.NotGuest, middleware.CanEditEvent).Patch("/", PatchEvent)
		r.With(middleware.NotGuest, middleware.CanEditEvent).Delete("/", DeleteEvent)

		r.Route("/positions", func(r chi.Router) {
			r.Get("/", GetEventPositions)
			r.With(middleware.NotGuest, middleware.CanEditEvent).Post("/", CreateEventPosition)

			r.Route("/{EventPositionID}", func(r chi.Router) {
				r.Use(positionCtx)

				r.Get("/", GetEventPosition)
				r.With(middleware.NotGuest, middleware.CanEditEvent).Patch("/", PatchEventPosition)
				r.With(middleware.NotGuest, middleware.CanEditEvent).Delete("/", DeleteEventPosition)
			})
		})

		r.Route("/signups", func(r chi.Router) {
			r.Get("/", GetEventSignups)
			r.With(middleware.NotGuest, middleware.CanEventSignup).Post("/", CreateEventSignup)

			r.Route("/{EventSignupID}", func(r chi.Router) {
				r.Use(signupCtx)

				r.Get("/", GetEventSignup)
				r.With(middleware.NotGuest, middleware.CanDeleteEventSignup).Delete("/", DeleteEventSignup)
			})
		})

		r.Route("/routing", func(r chi.Router) {
			r.Get("/", GetEventRouting)
			r.With(middleware.NotGuest, middleware.CanEditEvent).Post("/", CreateEventRouting)

			r.Route("/{EventRoutingID}", func(r chi.Router) {
				r.Use(routingCtx)

				r.With(middleware.NotGuest, middleware.CanEditEvent).Patch("/", PatchEventRouting)
				r.With(middleware.NotGuest, middleware.CanEditEvent).Delete("/", DeleteEventRouting)
			})
		})
	})
}

func TemplateRouter(r chi.Router) {
	r.With(middleware.NotGuest).Get("/", GetEventTemplates)

	r.Route("/{EventTemplateID}", func(r chi.Router) {
		r.Use(templateCtx)

		r.With(middleware.NotGuest, middleware.CanEditEvent).Put("/", UpdateEventTemplate)
		r.With(middleware.NotGuest, middleware.CanEditEvent).Delete("/", DeleteEventTemplate)
	})
}

func eventCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "EventID")
		if eventID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		uintEventID, err := strconv.ParseUint(eventID, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		ev := &models.Event{ID: uint(uintEventID)}
		if err := ev.Get(); err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.EventKey{}, ev)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func positionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		positionID := chi.URLParam(r, "EventPositionID")
		if positionID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		uintPositionID, err := strconv.ParseUint(positionID, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		ep := &models.EventPosition{ID: uint(uintPositionID)}
		if err := ep.Get(); err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.EventPositionKey{}, ep)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func signupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signupID := chi.URLParam(r, "EventSignupID")
		if signupID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		uintSignupID, err := strconv.ParseUint(signupID, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		es := &models.EventSignup{ID: uint(uintSignupID)}
		if err := es.Get(); err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.EventSignupKey{}, es)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func routingCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routingID := chi.URLParam(r, "EventRoutingID")
		if routingID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		uintRoutingID, err := strconv.ParseUint(routingID, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		er := &models.EventRouting{ID: uint(uintRoutingID)}
		if err := er.Get(); err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.EventRoutingKey{}, er)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func templateCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateID := chi.URLParam(r, "EventTemplateID")
		if templateID == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		uintTemplateID, err := strconv.ParseUint(templateID, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		et := &models.EventTemplate{ID: uint(uintTemplateID)}
		if err := et.Get(); err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.EventTemplateKey{}, et)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
