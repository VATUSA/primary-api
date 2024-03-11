package notification

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListNotifications)
	r.Post("/", CreateNotification)
	r.Route("/{NotificationID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Get("/", GetNotification)
		r.Put("/", UpdateNotification)
		r.Patch("/", PatchNotification)
		r.Delete("/", DeleteNotification)
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
