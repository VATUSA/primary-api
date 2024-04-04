package internal

import (
	_ "github.com/VATUSA/primary-api/internal/docs"
	v1 "github.com/VATUSA/primary-api/internal/v1"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           VATUSA Internal
// @version         0.1
// @description     VATUSAs internal API for use by other internal VATUSA services.
// @termsOfService  http://swagger.io/terms/

// @contact.name   VATUSA Support
// @contact.url    http://www.swagger.io/support
// @contact.email  vatusa6@vatusa.net

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.vatusa.dev
// @BasePath  /internal/v1
// @schemes https

func Router(r chi.Router, cfg *config.Config) {
	r.Route("/internal", func(r chi.Router) {
		v1.Router(r, cfg)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("https://api.vatusa.dev/internal/swagger/doc.json"),
		))
	})
}
