package user

import (
	"context"
	"github.com/VATUSA/primary-api/external/v3/roster"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.With(middleware.NotGuest).Get("/", GetSelf)

	r.Route("/{CID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.CanViewUser).Get("/", GetUser)
		r.With(middleware.CanEditUser).Put("/", UpdateUser)
		r.With(middleware.CanEditUser).Patch("/", PatchUser)

		r.With(middleware.CanViewUser).Get("/roster", roster.GetUserRosters)
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
