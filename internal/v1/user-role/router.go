package user_role

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Get("/", ListUserRoles)
	r.Post("/", CreateUserRoles)
	r.Route("/{UserRoleID}", func(r chi.Router) {
		r.Use(Ctx)
		r.Get("/", GetUserRole)
		r.Delete("/", DeleteUserRole)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "UserRoleID")
		if id == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		UserRoleID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userRole := &models.UserRole{ID: uint(UserRoleID)}
		if err = userRole.Get(); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserRoleKey{}, userRole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
