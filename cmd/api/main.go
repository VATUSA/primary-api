package main

import (
	"fmt"
	"github.com/VATUSA/primary-api/internal"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/models"
	gochi "github.com/VATUSA/primary-api/pkg/go-chi"
	"github.com/VATUSA/primary-api/pkg/storage"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	_ = godotenv.Load(".env")
	cfg := config.New()

	bucket, err := storage.NewS3Client(cfg.S3)
	if err != nil {
		panic(err)
	}

	storage.PublicBucket = bucket
	database.DB = database.Connect(cfg.Database)
	models.AutoMigrate()

	r := gochi.New(cfg)
	internal.Router(r, cfg)
	log.Fatalf("Err starting http server: %s", http.ListenAndServe(fmt.Sprintf(":%s", cfg.API.Port), r))
}
