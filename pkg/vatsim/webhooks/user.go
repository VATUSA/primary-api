package vatsim_webhooks

import (
	"encoding/json"
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
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

// TODO - does user creation webhook only send to VATUSA if they select division_id = USA?
// TODO - Rating changes? if VATUSA is the authority its not gonna send the webhook back to VATUSA correct? or it will?
func ConsumeMessage(rawMessage string) error {
	msg := Message{}
	if err := json.Unmarshal([]byte(rawMessage), &msg); err != nil {
		return err
	}

	for _, action := range msg.Actions {
		switch action.Type {
		case UserCreation:
			//
		case UserChanged:
			//
		default:
			// TODO - log unknown actionType

		}
	}

	return nil
}

func HandleUserCreation(action Action) error {
	return nil
}

func HandleUserChange(action Action) error {
	return nil
}
