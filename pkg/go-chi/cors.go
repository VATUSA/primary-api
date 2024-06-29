package go_chi

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/cors"
	"strings"
)

func NewCors(cfg *config.Config) cors.Options {
	allowedOrigins := strings.Split(cfg.Cors.AllowedOrigins, ",")
	return cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "x-guest", "x-user", "x-api-key"},
		AllowCredentials: true,
		MaxAge:           300,
	}
}
