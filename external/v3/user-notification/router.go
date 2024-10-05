package user_notification

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
	r.With(middleware.CanEditUser).Get("/", GetNotificationSettings)
	r.With(middleware.CanEditUser).Put("/", UpdateNotificationSettings)
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserCtx(r)

		userNotification := &models.UserNotification{CID: user.CID}
		if err := userNotification.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserNotificationKey{}, userNotification)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
