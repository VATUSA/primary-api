package user

import (
	"encoding/json"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/cookie"
	logger "github.com/VATUSA/primary-api/pkg/logging"
	"github.com/VATUSA/primary-api/pkg/oauth"
	"github.com/VATUSA/primary-api/pkg/utils"
	gonanoid "github.com/matoous/go-nanoid"
	"io"
	"net/http"
	"time"
)

type VATSIMUser struct {
	CID      string `json:"cid"`
	Personal struct {
		FirstName string `json:"name_first"`
		LastName  string `json:"name_last"`
		FullName  string `json:"name_full"`
		Email     string `json:"email"`
	} `json:"personal"`
	VATSIM struct {
		ControllerRating struct {
			ID    int    `json:"id"`
			Short string `json:"short"`
			Long  string `json:"long"`
		} `json:"rating"`
		PilotRating struct {
			ID    int    `json:"id"`
			Short string `json:"short"`
			Long  string `json:"long"`
		} `json:"pilotrating"`
		Region struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"region"`
		Division struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"division"`
		Subdivision struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"subdivision"`
	}
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
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

	utils.TempRedirect(w, r, oauth.OAuthConfig.AuthCodeURL(state))
}

func GetLoginCallback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		logger.ErrorWithErr(err, "Error getting session cookie")
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	session := make(map[string]string)
	if err := cookie.CookieStore.Decode("session", sessionCookie.Value, &session); err != nil {
		logger.ErrorWithErr(err, "Error decoding session cookie")
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	if r.URL.Query().Get("state") != session["state"] {
		logger.Error("State mismatch on login/callback")
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	token, err := oauth.OAuthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
	res, err := http.NewRequest("GET", config.Cfg.OAuth.UserInfoURL, nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	res.Header.Add("Accept", "application/json")
	res.Header.Add("User-Agent", "usa-primary-api")
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if resp.StatusCode >= 299 {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	user := &VATSIMUser{}
	if err := json.Unmarshal(body, user); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if encoded, err := cookie.CookieStore.Encode("VATUSA", map[string]string{
		"cid": user.CID,
	}); err == nil {
		cookie := &http.Cookie{
			Name:     "VATUSA",
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Domain:   config.Cfg.Cookie.Domain,
		}

		http.SetCookie(w, cookie)
	} else {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	redirect := session["redirect"]
	if redirect == "" {
		redirect = "https://vatusa.net"
	}

	// Delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   config.Cfg.Cookie.Domain,
		Expires:  time.Unix(0, 0),
	})

	utils.TempRedirect(w, r, redirect)
}

func GetLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "VATUSA",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   config.Cfg.Cookie.Domain,
		Expires:  time.Unix(0, 0),
	})

	utils.Response(r, http.StatusNoContent)
}
