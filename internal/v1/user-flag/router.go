package user_flag

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(r chi.Router) {
	r.Use(Ctx)
	r.Get("/", GetUserFlag)
	r.Put("/", UpdateUserFlag)
	r.Delete("/", DeleteUserFlag)
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserCtx(r)

		userFlag := &models.UserFlag{CID: user.CID}
		if err := userFlag.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "userFlag", userFlag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
