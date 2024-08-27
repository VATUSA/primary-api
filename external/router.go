package external

import (
	"fmt"
	"github.com/VATUSA/primary-api/external/docs"
	_ "github.com/VATUSA/primary-api/external/docs"
	v3 "github.com/VATUSA/primary-api/external/v3"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"strings"
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

// @BasePath  /v3

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key

func Router(r chi.Router, cfg *config.Config) {
	v3.Router(r, cfg)

	docs.SwaggerInfo.Host = cfg.API.BaseURL[strings.Index(cfg.API.BaseURL, "://")+3:]

	if strings.Contains(cfg.API.BaseURL, "https") {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", config.Cfg.API.BaseURL)),
	))
}
