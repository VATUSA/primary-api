package notification

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.With(middleware.CanViewNotifications).Get("/", ListNotifications)

	r.Route("/{NotificationID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.CanEditNotifications).Put("/", UpdateNotification)
		r.With(middleware.CanEditNotifications).Patch("/", PatchNotification)
		r.With(middleware.CanEditNotifications).Delete("/", DeleteNotification)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "NotificationID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		NotificationID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		notification := &models.Notification{ID: uint(NotificationID)}
		if err = notification.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.NotificationKey{}, notification)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
