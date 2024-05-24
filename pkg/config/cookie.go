package config

type Cookie struct {
	HashKey  []byte
	BlockKey []byte
	Domain   string
}

func NewCookie() *Cookie {
	return &Cookie{
		HashKey:  []byte(EnvOrDefault("COOKIE_HASH_KEY", string(defaultCfg.Cookie.HashKey))),
		BlockKey: []byte(EnvOrDefault("COOKIE_BLOCK_KEY", string(defaultCfg.Cookie.BlockKey))),
		Domain:   EnvOrDefault("COOKIE_DOMAIN", defaultCfg.Cookie.Domain),
	}
}
