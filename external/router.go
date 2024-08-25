package external

import (
	_ "github.com/VATUSA/primary-api/external/docs"
	v3 "github.com/VATUSA/primary-api/external/v3"
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
// @BasePath  /v3
// @schemes https

func Router(r chi.Router, cfg *config.Config) {
	v3.Router(r, cfg)

	r.Get("/swagger/*", httpSwagger.Handler(
		//httpSwagger.URL("https://api.vatusa.dev/swagger/doc.json"),
		httpSwagger.URL("http://api.vatusa.local:3000/swagger/doc.json"),
	))

}
