package main

import (
	"fmt"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/cookie"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/models"
	gochi "github.com/VATUSA/primary-api/pkg/go-chi"
	"github.com/VATUSA/primary-api/pkg/oauth"
	"github.com/VATUSA/primary-api/pkg/storage"
	"github.com/VATUSA/primary-api/views"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	_ = godotenv.Load(".env")
	config.Cfg = config.New()

	log.SetLevel(log.DebugLevel)

	oauth.OAuthConfig = oauth.InitializeVATSIM(config.Cfg)
	oauth.DiscordOAuthConfig = oauth.InitializeDiscord(config.Cfg)

	bucket, err := storage.NewS3Client(config.Cfg.S3)
	if err != nil {
		panic(err)
	}

	storage.PublicBucket = bucket
	database.DB = database.Connect(config.Cfg.Database)
	cookie.CookieStore = cookie.New(config.Cfg)
	models.AutoMigrate()

	r := gochi.New(config.Cfg)
	views.Router(r, config.Cfg)
	log.Fatalf("Err starting http server: %s", http.ListenAndServe(fmt.Sprintf(":%s", config.Cfg.API.Port), r))
}
