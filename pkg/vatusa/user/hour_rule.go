package user

import (
	vatsim_api "github.com/VATUSA/primary-api/pkg/vatsim/api"
	"github.com/VATUSA/primary-api/pkg/vnas"
	log "github.com/sirupsen/logrus"
	"time"
)

func FiftyFiftyRuleCheck(homeARTCC string, cid uint) (bool, error) {
	artcc, err := vnas.GetARTCC(homeARTCC)
	if err != nil {
		return false, err
	}

	positions := artcc.Positions()
	positionsMapped := make(map[string]bool, len(positions))
	for _, pos := range positions {
		positionsMapped[pos] = true
	}

	connections, err := vatsim_api.GetATCConnections(cid)
	if err != nil {
		return false, err
	}

	var homeHours time.Duration
	var otherHours time.Duration

	for _, connection := range connections.Items {
		if time.Since(connection.ConnectionId.End) > 30*24*time.Hour {
			break
		}
		hours := connection.ConnectionId.End.Sub(connection.ConnectionId.Start)

		if positionsMapped[connection.ConnectionId.Callsign] {
			homeHours += hours
		} else {
			otherHours += hours
		}
	}

	log.Debug("Home hours: ", homeHours)
	log.Debug("Other hours: ", otherHours)

	return homeHours > otherHours, nil
}
