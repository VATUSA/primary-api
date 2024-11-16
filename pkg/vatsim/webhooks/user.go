package vatsim_webhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Message struct {
	CID     uint     `json:"resource"`
	Actions []Action `json:"actions"`
}

type WebhookActionType string

const (
	UserCreation WebhookActionType = "member_created_action"
	UserChanged  WebhookActionType = "member_changed_action"
)

type Action struct {
	Type      WebhookActionType `json:"action"`
	Authority string            `json:"authority"`
	Comment   string            `json:"comment"`
	Deltas    []Delta           `json:"deltas"`
	Timestamp time.Time         `json:"timestamp"`
}

type Delta struct {
	Field  Field       `json:"field"`
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

func ConsumeMessage(rawMessage string) error {
	msg := Message{}
	if err := json.Unmarshal([]byte(rawMessage), &msg); err != nil {
		return err
	}

	for _, action := range msg.Actions {
		switch action.Type {
		case UserCreation:
			if action.Authority != "myVATSIM" {
				log.Warnf("Unknown authority: %s", action.Authority)
				return errors.New("unknown authority")
			}

			return handleUserCreation(action)
		case UserChanged:
			return handleUserChange(msg.CID, action)
		default:
			log.Errorf("Unknown action type: %s", action.Type)
			return errors.New("unknown action type")
		}
	}

	return nil
}

func handleUserCreation(action Action) error {
	newUser := &models.User{
		LastLogin:    time.Time{},
		LastCertSync: time.Now(),
		UpdatedAt:    time.Now(),
	}

	for _, delta := range action.Deltas {
		switch delta.Field {
		case FieldID:
			newUser.CID = delta.After.(uint)
		case FieldNameFirst:
			newUser.FirstName = delta.After.(string)
		case FieldNameLast:
			newUser.LastName = delta.After.(string)
		case FieldEmail:
			newUser.Email = delta.After.(string)
		case FieldRating:
			newUser.ControllerRating = constants.ATCRating(delta.After.(int))
		case FieldPilotRating:
			newUser.PilotRating = constants.PilotRating(delta.After.(int))
		case FieldRegistrationDate:
			newUser.CreatedAt = delta.After.(time.Time)
		}
	}

	newUser.PreferredOIs = fmt.Sprintf("%s%s", strings.ToUpper(newUser.FirstName[:1]), strings.ToUpper(newUser.LastName[:1]))

	if err := newUser.Create(); err != nil {
		fmt.Printf("Error creating user: %s\n", err.Error())
		return err
	}

	userFlag := &models.UserFlag{
		CID: newUser.CID,
	}

	if err := userFlag.Create(); err != nil {
		fmt.Printf("Error creating user flag: %s\n", err.Error())
		return err
	}

	// Create User Notification Settings
	userNotificationSettings := &models.UserNotification{
		CID:            newUser.CID,
		DiscordEnabled: true,
		EmailEnabled:   true,
		Events:         true,
		Training:       true,
		Feedback:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := userNotificationSettings.Create(); err != nil {
		fmt.Printf("Error creating user notification settings: %s\n", err.Error())
		return err
	}

	return nil
}

func handleUserChange(cid uint, action Action) error {
	user := &models.User{
		CID: cid,
	}

	if err := user.Get(); err != nil {
		return err
	}

	for _, delta := range action.Deltas {
		switch delta.Field {
		case FieldID:
			user.CID = delta.After.(uint)
		case FieldEmail:
			user.Email = delta.After.(string)
		case FieldRating:
			oldRating := user.ControllerRating
			user.ControllerRating = constants.ATCRating(delta.After.(int))
			rc := &models.RatingChange{
				CID:          user.CID,
				OldRating:    oldRating,
				NewRating:    user.ControllerRating,
				CreatedAt:    time.Now(),
				CreatedByCID: 0,
				UpdatedAt:    time.Now(),
			}

			if err := rc.Create(); err != nil {
				return err
			}
		case FieldPilotRating:
			user.PilotRating = constants.PilotRating(delta.After.(int))
		}
	}

	if err := user.Update(); err != nil {
		return err
	}

	if user.ControllerRating == constants.SuspendedRating {
		// Create a disciplinary log entry
		disciplinaryLogEntry := &models.DisciplinaryLogEntry{
			CID:        user.CID,
			Entry:      "User controller rating changed to suspended(0).",
			VATUSAOnly: false,
			CreatedAt:  time.Now(),
			CreatedBy:  "VATSIM Webhook",
			UpdatedAt:  time.Now(),
			UpdatedBy:  "",
		}

		if err := disciplinaryLogEntry.Create(); err != nil {
			return err
		}

		userFlag := &models.UserFlag{
			CID: user.CID,
		}

		if err := userFlag.Get(); err != nil {
			return err
		}

		userFlag.NoStaffRole = true
		userFlag.NoStaffLogEntryID = disciplinaryLogEntry.ID
		if err := userFlag.Update(); err != nil {
			return err
		}
	}

	return nil
}
