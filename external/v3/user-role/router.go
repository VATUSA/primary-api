package user_role

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router(r chi.Router) {
	r.Get("/", GetSelfRoles)

	r.Route("/{RoleID}", func(r chi.Router) {
		r.Use(Ctx)

		r.With(middleware.CanAddRole).Get("/", CreateUserRoles)
		r.With(middleware.CanDeleteRole).Delete("/", DeleteUserRoles)
	})
}

func Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleId := chi.URLParam(r, "RoleID")
		if roleId == "" {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		role := &models.UserRole{RoleID: constants.RoleID(roleId)}
		err := role.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserRoleKey{}, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
