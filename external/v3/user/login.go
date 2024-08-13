package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/cookie"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/oauth"
	"github.com/VATUSA/primary-api/pkg/utils"
	vatsim_api "github.com/VATUSA/primary-api/pkg/vatsim/api"
	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	session := make(map[string]string)
	if err := cookie.CookieStore.Decode("session", sessionCookie.Value, &session); err != nil {
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	if r.URL.Query().Get("state") != session["state"] {
		utils.Render(w, r, utils.ErrForbidden)
		return
	}

	token, err := exchangeToken(r.Context(), oauth.OAuthConfig, r.URL.Query().Get("code"))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	//token, err := oauth.OAuthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err)
	//	utils.Render(w, r, utils.ErrInternalServer)
	//	return
	//}

	res, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.Cfg.OAuth.BaseURL, config.Cfg.OAuth.UserInfoURL), nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	res.Header.Add("Accept", "application/json")
	res.Header.Add("User-Agent", "vatusa-primary-api")
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
		fmt.Println("Error Code:", resp.StatusCode)
		fmt.Println("Error Body:", string(body))
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	user := &vatsim_api.User{}
	if err := json.Unmarshal(body, user); err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
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
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	// Create user if they don't exist
	intCID, err := strconv.ParseInt(user.CID, 10, 64)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}
	dbUser := &models.User{
		CID: uint(intCID),
	}

	if err := dbUser.Get(); err != nil {
		dbUser.FirstName = user.Personal.FirstName
		dbUser.LastName = user.Personal.LastName
		dbUser.Email = user.Personal.Email
		dbUser.PilotRating = constants.PilotRating(user.VATSIM.PilotRating.ID)
		dbUser.ControllerRating = constants.ATCRating(user.VATSIM.ControllerRating.ID)
		dbUser.LastLogin = time.Now()
		dbUser.LastCertSync = time.Now()
		if err := dbUser.Create(); err != nil {
			utils.Render(w, r, utils.ErrInternalServerWithErr(err))
			return
		}

		// Create user flags
		dbUserFlag := &models.UserFlag{
			CID:                  dbUser.CID,
			NoStaffRole:          false,
			NoVisiting:           false,
			NoTransferring:       false,
			NoTraining:           false,
			UsedTransferOverride: false,
		}

		if err := dbUserFlag.Create(); err != nil {
			utils.Render(w, r, utils.ErrInternalServerWithErr(err))
			return
		}
	} else {
		dbUser.FirstName = user.Personal.FirstName
		dbUser.LastName = user.Personal.LastName
		dbUser.Email = user.Personal.Email
		dbUser.PilotRating = constants.PilotRating(user.VATSIM.PilotRating.ID)
		dbUser.ControllerRating = constants.ATCRating(user.VATSIM.ControllerRating.ID)
		dbUser.LastLogin = time.Now()
		dbUser.LastCertSync = time.Now()
		if err := dbUser.Update(); err != nil {
			utils.Render(w, r, utils.ErrInternalServerWithErr(err))
			return
		}

		// Create user flags if they don't exist (due to old data migration)
		dbUserFlag := &models.UserFlag{
			CID: dbUser.CID,
		}

		if err := dbUserFlag.Get(); err != nil {
			dbUserFlag.NoStaffRole = false
			dbUserFlag.NoVisiting = false
			dbUserFlag.NoTransferring = false
			dbUserFlag.NoTraining = false
			dbUserFlag.UsedTransferOverride = false
			if err := dbUserFlag.Create(); err != nil {
				utils.Render(w, r, utils.ErrInternalServerWithErr(err))
				return
			}
		}
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

func exchangeToken(ctx context.Context, oauthConfig *oauth2.Config, code string) (*oauth2.Token, error) {
	tokenURL := oauthConfig.Endpoint.TokenURL

	// Prepare form data
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {oauthConfig.RedirectURL},
	}

	if oauthConfig.ClientSecret != "" {
		data.Set("client_id", oauthConfig.ClientID)
		data.Set("client_secret", oauthConfig.ClientSecret)
	}

	reqBody := strings.NewReader(data.Encode())

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request (using the default client or a custom one)
	client := http.DefaultClient // Or use your custom client
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected token response status: %s", resp.Status)
	}

	// Parse the response (adjust for your actual response structure)
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding token response: %w", err)
	}

	// Create the OAuth2 token
	token := &oauth2.Token{
		AccessToken: tokenResponse.AccessToken,
		TokenType:   tokenResponse.TokenType,
		Expiry:      time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
	}

	return token, nil
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
