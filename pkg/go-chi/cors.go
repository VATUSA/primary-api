package go_chi

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"strings"
)

func NewCors(cfg *config.Config) cors.Options {
	allowedOrigins := strings.Split(cfg.Cors.AllowedOrigins, ",")
	log.Debug("Allowed Origins: ", allowedOrigins)
	return cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "x-guest", "x-user", "x-api-key"},
		AllowCredentials: true,
		MaxAge:           300,
	}
}
