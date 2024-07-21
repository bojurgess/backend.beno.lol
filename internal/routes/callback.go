package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/spotify"
	"github.com/bojurgess/backend.beno.lol/internal/util"
)

type Callback struct {
	DB     *database.Database
	Config *config.Config
}

func (p *Callback) Route(w http.ResponseWriter, r *http.Request) {
	if *p.Config.Mode == "production" || *p.Config.Mode == "prod" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var storedState string
	requestState := r.URL.Query().Get("state")

	// handle state
	if cookie, err := r.Cookie("state"); err == nil {
		storedState = cookie.Value

		if storedState != requestState {
			http.Error(w, "state mismatch", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "state mismatch", http.StatusBadRequest)
		return
	}

	// handle error
	if err := r.URL.Query().Get("error"); err != "" {
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	var protocol string

	if *p.Config.Origin != "backend.beno.lol" {
		protocol = "http"
	} else {
		protocol = "https"
	}

	code := r.URL.Query().Get("code")

	url := "https://accounts.spotify.com/api/token"
	body := util.MapToQuerystring(map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": fmt.Sprintf("%s://%s/auth/callback", protocol, *p.Config.Origin),
	})
	headers := map[string]string{
		"Authorization": util.EncodeAuth(p.Config.Env.SpotifyClientID, p.Config.Env.SpotifyClientSecret),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tr spotify.TokenResponse
	if err = json.NewDecoder(res.Body).Decode(&tr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tr.TokenErrorResponse != nil {
		http.Error(w, tr.Error, http.StatusBadRequest)
		return
	}

	tokens := spotify.MapTokenResponse(tr)

	user, err := spotify.GetUser(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := p.DB.AddUser(*user); err != nil {
		// this is so ugly!! >:(
		if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = p.DB.UpdateUser(*user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Successfully authenticated! You can now close this window."))
}
