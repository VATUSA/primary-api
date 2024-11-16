package vatsim_webhooks

import (
	"errors"
	"net/http"
)

func GetAuthentication(r *http.Request) (string, error) {
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "VATSIM-API" {
		return "", errors.New("invalid user agent")
	}

	authHeader := r.Header.Get("Authentication")
	if authHeader == "" {
		return "", errors.New("missing authentication header")
	}

	return authHeader, nil
}
