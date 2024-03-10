package external

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           VATUSA API
// @version         0.1
// @description     VATUSAs public API
// @termsOfService  http://swagger.io/terms/

// @contact.name   VATUSA Support
// @contact.url    http://www.swagger.io/support
// @contact.email vatusa6@vatusa.net

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.vatusa.net
// @BasePath  /v1
// @schemes http

func Router(r chi.Router, cfg *config.Config) {
	r.Route("/v1", func(r chi.Router) {

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://api.vatusa.net/v1/swagger/doc.json"),
		))
	})
}
