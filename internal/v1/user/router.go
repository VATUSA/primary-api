package user

import (
	"context"
	action_log "github.com/VATUSA/primary-api/internal/v1/action-log"
	disciplinary_log "github.com/VATUSA/primary-api/internal/v1/disciplinary-log"
	rating_change "github.com/VATUSA/primary-api/internal/v1/rating-change"
	"github.com/VATUSA/primary-api/internal/v1/roster"
	roster_request "github.com/VATUSA/primary-api/internal/v1/roster-request"
	user_flag "github.com/VATUSA/primary-api/internal/v1/user-flag"
	user_role "github.com/VATUSA/primary-api/internal/v1/user-role"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListUsers)
	r.Post("/", CreateUser)

	r.Route("/{CID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Get("/", GetUser)
		r.Put("/", UpdateUser)
		r.Patch("/", PatchUser)
		r.Delete("/", DeleteUser)

		r.Route("/action-log", func(r chi.Router) {
			action_log.Router(r)
		})

		r.Route("/disciplinary-log", func(r chi.Router) {
			disciplinary_log.Router(r)
		})

		r.Route("/rating-change", func(r chi.Router) {
			rating_change.Router(r)
		})

		r.Route("/roster", func(r chi.Router) {
			roster.Router(r)
		})

		r.Route("/roster-request", func(r chi.Router) {
			roster_request.Router(r)
		})

		r.Route("/user-flag", func(r chi.Router) {
			user_flag.Router(r)
		})

		r.Route("/role", func(r chi.Router) {
			user_role.Router(r)
		})

	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := chi.URLParam(r, "CID")
		if cid == "" {
			render.Render(w, r, utils.ErrNotFound)
			return
		}

		CID, err := strconv.ParseUint(cid, 10, 64)
		if err != nil {
			render.Render(w, r, utils.ErrNotFound)
			return
		}

		user := &models.User{CID: uint(CID)}
		err = user.Get()
		if err != nil {
			render.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
