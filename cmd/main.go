package main

import (
	"encoding/json"
	"github.com/VATUSA/primary-api/internal/training"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

	//check, err := user.FiftyFiftyRuleCheck("ZDV", 811918)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Info("FiftyFiftyRuleCheck: ", check)
	_ = godotenv.Load(".env")
	config.Cfg = config.New()
	database.DB = database.Connect(config.Cfg.Database)

	s2OTS := training.S2OTS
	jsonS2OTS, err := json.Marshal(s2OTS)
	if err != nil {
		log.Fatal(err)
	}

	otsTemplate := &models.OTSTemplate{
		Name:          "S2 OTS",
		Template:      jsonS2OTS,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastUpdatedBy: 1293257,
	}

	if err := otsTemplate.Create(); err != nil {
		log.Fatal(err)
	}
}
