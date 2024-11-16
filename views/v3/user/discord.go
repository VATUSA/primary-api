package user

import (
	"encoding/json"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/cookie"
	"github.com/VATUSA/primary-api/pkg/oauth"
	"github.com/VATUSA/primary-api/pkg/utils"
	gonanoid "github.com/matoous/go-nanoid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func GetDiscordLink(w http.ResponseWriter, r *http.Request) {
	state, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	sessionEncoded, err := cookie.CookieStore.Encode("session", map[string]string{
		"state":    state,
		"redirect": r.URL.Query().Get("redirect"),
	})
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionEncoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   config.Cfg.Cookie.Domain,
	}

	http.SetCookie(w, cookie)

	utils.TempRedirect(w, r, oauth.DiscordOAuthConfig.AuthCodeURL(state))
}

type DiscordResp struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func GetDiscordCallback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.WithError(err).Error("Error getting session cookie.")
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	session := make(map[string]string)
	if err := cookie.CookieStore.Decode("session", sessionCookie.Value, &session); err != nil {
		log.WithError(err).Error("Error decoding session cookie.")
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	if r.URL.Query().Get("state") != session["state"] {
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	token, err := exchangeToken(r.Context(), oauth.DiscordOAuthConfig, r.URL.Query().Get("code"))
	if err != nil {
		log.WithError(err).Errorf("Error exchanging tokens with Discord, used code: %s.", r.URL.Query().Get("code"))
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	//token, err := oauth.OAuthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err)
	//	utils.Render(w, r, utils.ErrInternalServer)
	//	return
	//}

	res, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.Cfg.DiscordOAuth.BaseURL, config.Cfg.DiscordOAuth.UserInfoURL), nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	res.Header.Add("Accept", "application/json")
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	if resp.StatusCode >= 299 {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	discordResp := &DiscordResp{}
	if err := json.Unmarshal(body, discordResp); err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	user := utils.GetXUser(r)
	user.DiscordID = discordResp.ID
	if err := user.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	redirect := session["redirect"]
	if redirect == "" {
		redirect = "https://vatusa.net"
	} else {
		redirect = fmt.Sprintf("%s?name=%s#discord", redirect, discordResp.Username)
	}

	utils.TempRedirect(w, r, redirect)
}

// UnlinkDiscord godoc
// @Summary Unlink your Discord account
// @Description Unlink your Discord account from your VATUSA account
// @Tags discord
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/discord/unlink [get]
func UnlinkDiscord(w http.ResponseWriter, r *http.Request) {
	user := utils.GetXUser(r)

	if user.DiscordID == "" {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	user.DiscordID = ""
	if err := user.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusOK)
}
