package main

//func main() {
//	_ = godotenv.Load(".env")
//	config.Cfg = config.New()
//
//	oauth.OAuthConfig = oauth.Initialize(config.Cfg)
//
//	bucket, err := storage.NewS3Client(config.Cfg.S3)
//	if err != nil {
//		panic(err)
//	}
//
//	storage.PublicBucket = bucket
//	database.DB = database.Connect(config.Cfg.Database)
//	cookie.CookieStore = cookie.New(config.Cfg)
//	models.AutoMigrate()
//	logger.Setup()
//
//	r := gochi.New(config.Cfg)
//	external.Router(r, config.Cfg)
//	log.Fatalf("Err starting http server: %s", http.ListenAndServeTLS(fmt.Sprintf(":%s", config.Cfg.API.Port), "./api.vatusa.local.crt", "./api.vatusa.local.key", r))
//}
