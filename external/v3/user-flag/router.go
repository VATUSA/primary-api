package user_flag

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(r chi.Router) {
	r.Use(middleware.NotGuest, Ctx)
	r.With(middleware.CanViewUserFlag).Get("/", GetUserFlag)
	r.With(middleware.CanEditUserFlag).Put("/", UpdateUserFlag)
	r.With(middleware.CanEditUserFlag).Patch("/", PatchUserFlag)
	r.With(middleware.CanEditUserFlag).Delete("/", DeleteUserFlag)
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserCtx(r)

		userFlag := &models.UserFlag{CID: user.CID}
		if err := userFlag.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserFlagKey{}, userFlag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
