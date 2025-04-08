package vatsim_api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	baseURL   = "https://api.vatsim.net/api"
	baseV2URL = "https://api.vatsim.net/v2"
	dataURL   = "https://data.vatsim.net/v3/vatsim-data.json"
)

type Location struct {
	Region      string `json:"region"`
	Division    string `json:"division"`
	Subdivision string `json:"subdivision"`
}

func GetLocation(cid uint) (Location, error) {
	resp, err := http.Get(fmt.Sprintf("%s/ratings/%d/", baseURL, cid))
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}

	if resp.StatusCode > 299 {
		log.Warnf("Failed to get division for %s: %s", cid, body)
		return Location{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var division Location
	err = json.Unmarshal(body, &division)
	if err != nil {
		return Location{}, err
	}

	return division, nil
}

type UserStats struct {
	Atc   float32 `json:"atc"`
	Pilot float32 `json:"pilot"`
	S1    float32 `json:"s1"`
	S2    float32 `json:"s2"`
	S3    float32 `json:"s3"`
	C1    float32 `json:"c1"`
	C2    float32 `json:"c2"`
	C3    float32 `json:"c3"`
	I1    float32 `json:"i1"`
	I2    float32 `json:"i2"`
	I3    float32 `json:"i3"`
	Sup   float32 `json:"sup"`
	Adm   float32 `json:"adm"`
}

func GetHours(cid uint) (UserStats, error) {
	resp, err := http.Get(fmt.Sprintf("%s/members/%d/stats", baseV2URL, cid))
	if err != nil {
		return UserStats{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserStats{}, err
	}

	if resp.StatusCode > 299 {
		log.Warnf("Failed to get stats for %s: %s", cid, body)
		return UserStats{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var userStats UserStats
	err = json.Unmarshal(body, &userStats)
	if err != nil {
		return UserStats{}, err
	}

	return userStats, nil
}

type ATCConnections struct {
	Items []struct {
		ConnectionId struct {
			Id       int       `json:"id"`
			VatsimId string    `json:"vatsim_id"`
			Type     int       `json:"type"`
			Rating   int       `json:"rating"`
			Callsign string    `json:"callsign"`
			Start    time.Time `json:"start"`
			End      time.Time `json:"end"`
			Server   string    `json:"server"`
		} `json:"connection_id"`
		AircraftTracked int `json:"aircrafttracked"`
		AircraftSeen    int `json:"aircraftseen"`
	} `json:"items"`
	Count uint `json:"count"`
}

// GetATCConnections - for the last 30 days
func GetATCConnections(cid uint) (ATCConnections, error) {
	resp, err := http.Get(fmt.Sprintf("%s/members/%d/atc", baseV2URL, cid))
	if err != nil {
		return ATCConnections{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ATCConnections{}, err
	}

	if resp.StatusCode > 299 {
		log.Warnf("Failed to get stats for %s: %s", cid, body)
		return ATCConnections{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var atcConnections ATCConnections
	err = json.Unmarshal(body, &atcConnections)
	if err != nil {
		return ATCConnections{}, err
	}

	return atcConnections, nil
}
