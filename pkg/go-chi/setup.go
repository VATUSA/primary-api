package go_chi

import (
	"github.com/VATUSA/primary-api/pkg/config"
	auth "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"net/http"
)

func New(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(NewCors(cfg)))

	r.Use(auth.HasCookie)
	r.Use(auth.HasAPIKey)

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusNoContent)
		})
	})

	return r
}
