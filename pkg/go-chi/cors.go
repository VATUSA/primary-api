package go_chi

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/go-chi/cors"
)

func NewCors(cfg *config.Config) cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{cfg.Cors.AllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "x-guest", "x-user", "x-api-key"},
		AllowCredentials: true,
		MaxAge:           300,
	}
}
