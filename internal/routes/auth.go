package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/util"
)

type Auth struct {
	DB     *database.Database
	Config *config.Config
}

func (p *Auth) Route(w http.ResponseWriter, r *http.Request) {
	var protocol string

	if *p.Config.Mode == "production" || *p.Config.Mode == "prod" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if *p.Config.Origin != "backend.beno.lol" {
		protocol = "http"
	} else {
		protocol = "https"
	}

	config := p.Config

	state := randomString(15)
	scope := []string{
		"user-read-private",
		"user-read-email",
		"user-read-currently-playing",
	}

	fmt.Println(protocol)

	url := "https://accounts.spotify.com/authorize?"
	query := util.MapToQuerystring(map[string]string{
		"client_id":     config.Env.SpotifyClientID,
		"response_type": "code",
		"redirect_uri":  fmt.Sprintf("%s://%s/auth/callback", protocol, *config.Origin),
		"state":         state,
		"scope":         strings.Join(scope, " "),
		"show_dialog":   "true",
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "state",
		Value: state,

		// 5 minute expiry
		MaxAge: 60 * 5,

		HttpOnly: true,
		Secure:   *config.Mode == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, url+query, http.StatusMovedPermanently)
}

func randomString(len int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := strings.Builder{}

	for i := 0; i < len; i++ {
		b.WriteByte(charset[rand.Intn(len)])
	}

	return b.String()
}
