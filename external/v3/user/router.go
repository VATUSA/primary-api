package user

import (
	"context"
	action_log "github.com/VATUSA/primary-api/external/v3/action-log"
	disciplinary_log "github.com/VATUSA/primary-api/external/v3/disciplinary-log"
	"github.com/VATUSA/primary-api/external/v3/feedback"
	"github.com/VATUSA/primary-api/external/v3/notification"
	rating_change "github.com/VATUSA/primary-api/external/v3/rating-change"
	"github.com/VATUSA/primary-api/external/v3/roster"
	user_flag "github.com/VATUSA/primary-api/external/v3/user-flag"
	user_notification "github.com/VATUSA/primary-api/external/v3/user-notification"
	user_role "github.com/VATUSA/primary-api/external/v3/user-role"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.With(middleware.NotGuest).Get("/logout", GetLogout)

	r.Get("/login", GetLogin)
	r.Get("/login/callback", GetLoginCallback)

	r.With(middleware.NotGuest).Get("/discord", GetDiscordLink)
	r.With(middleware.NotGuest).Get("/discord/callback", GetDiscordCallback)
	r.With(middleware.NotGuest).Get("/discord/unlink", UnlinkDiscord)

	r.With(middleware.NotGuest).Get("/", GetSelf)

	r.Route("/{CID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.NotGuest, middleware.CanViewUser).Get("/", GetUser)
		r.With(middleware.NotGuest, middleware.CanEditUser).Put("/", UpdateUser)
		r.With(middleware.NotGuest, middleware.CanEditUser).Patch("/", PatchUser)

		r.Route("/action-log", func(r chi.Router) {
			action_log.Router(r)
		})

		r.Route("/disciplinary-log", func(r chi.Router) {
			disciplinary_log.Router(r)
		})

		r.With(middleware.NotGuest, middleware.CanViewUser).Get("/feedback", feedback.GetUserFeedback)

		r.Route("/notifications", func(r chi.Router) {
			notification.Router(r)
		})

		r.Route("/notification-settings", func(r chi.Router) {
			user_notification.Router(r)
		})

		r.Route("/rating-change", func(r chi.Router) {
			rating_change.Router(r)
		})

		r.With(middleware.NotGuest, middleware.CanViewUser).Get("/roster", roster.GetUserRosters)

		r.Route("/user-flag", func(r chi.Router) {
			user_flag.Router(r)
		})

		r.Route("/roles", func(r chi.Router) {
			user_role.Router(r)
		})
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := chi.URLParam(r, "CID")
		if cid == "" {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		CID, err := strconv.ParseUint(cid, 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		user := &models.User{CID: uint(CID)}
		err = user.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
