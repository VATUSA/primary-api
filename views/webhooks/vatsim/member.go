package vatsim

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	vatsim_webhooks "github.com/VATUSA/primary-api/pkg/vatsim/webhooks"
	"io"
	"net/http"
)

func ProcessMemberWebhook(w http.ResponseWriter, r *http.Request) {
	authToken, err := vatsim_webhooks.GetAuthentication(r)
	if err != nil {
		utils.Render(w, r, utils.ErrUnauthorized)
		return
	}

	if authToken != "" {
		utils.Render(w, r, utils.ErrUnauthorized)
		return
	}

	// Parse body into string
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
	body := string(bodyBytes)

	// Process the webhook
	if err := vatsim_webhooks.ConsumeMessage(body); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusOK)
}
