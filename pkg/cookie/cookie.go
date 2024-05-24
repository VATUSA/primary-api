package cookie

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/gorilla/securecookie"
)

var CookieStore *securecookie.SecureCookie

func New(cfg *config.Config) *securecookie.SecureCookie {
	return securecookie.New(cfg.Cookie.HashKey, cfg.Cookie.BlockKey)
}
