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
	config := p.Config

	state := randomString(15)
	scope := []string{
		"user-read-private",
		"user-read-email",
		"user-read-currently-playing",
	}

	url := "https://accounts.spotify.com/authorize?"
	query := util.MapToQuerystring(map[string]string{
		"client_id":     config.Env.SpotifyClientID,
		"response_type": "code",
		"redirect_uri":  fmt.Sprintf("http://%s:%d/auth/callback", *config.Host, *config.Port),
		"state":         state,
		"scope":         strings.Join(scope, " "),
		"show_dialog":   "true",
	})

	r.AddCookie(&http.Cookie{
		Name:  "state",
		Value: state,

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
