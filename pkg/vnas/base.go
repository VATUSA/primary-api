package vnas

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://data-api.vnas.vatsim.net/api"
)

type ARTCC struct {
	ID            string    `json:"id"`
	Facility      Facility  `json:"facility"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

type Facility struct {
	ID              string     `json:"id" example:"ZDV"`
	Name            string     `json:"name" example:"Denver ARTCC"`
	Positions       []Position `json:"positions"`
	ChildFacilities []Facility `json:"childFacilities"`
}

type Position struct {
	ID       string `json:"id" example:""`
	Name     string `json:"name" example:"14 - Hayden High"`
	Callsign string `json:"callsign" example:"DEN_14_CTR"`
}

func GetARTCC(id string) (*ARTCC, error) {
	resp, err := http.Get(fmt.Sprintf("%s/artccs/%s", baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Warnf("Failed to get ARTCC for %s: %s", id, body)
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var artcc ARTCC
	if err = json.Unmarshal(body, &artcc); err != nil {
		return nil, err
	}

	return &artcc, nil
}

func (artcc ARTCC) Positions() []string {
	var positions []string
	for _, position := range artcc.Facility.Positions {
		positions = append(positions, position.Callsign)
	}

	for _, facility := range artcc.Facility.ChildFacilities {
		for _, position := range facility.Positions {
			positions = append(positions, position.Callsign)
		}
	}

	return positions
}
