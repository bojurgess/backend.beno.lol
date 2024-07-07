package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/spotify"
)

type User struct {
	DB     *database.Database
	Config *config.Config
}

// Handles streaming of user NowPlaying data down to client.
// This is done through SSE.
// Route expects to be handled on a route with path parameter {id}.
func (p *User) Route(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	u, err := p.DB.GetUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// send data every second, until client disconnects
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Client disconnected")
			return
		case <-time.After(1 * time.Second):
			np, err := p.getNowPlaying(&u)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Fprintf(w, "data: %v\n\n", np)

			w.(http.Flusher).Flush()
		}
	}
}

// making a wrapper function for GetNowPlaying
func (p *User) getNowPlaying(u *database.User) (*spotify.NowPlayingResponse, error) {
	if time.Now().After(u.Tokens.ExpiresAt) {
		if err := spotify.RefreshAccessToken(&u.Tokens, p.Config.Env.SpotifyClientID, p.Config.Env.SpotifyClientSecret); err != nil {
			fmt.Printf("Error refreshing token: %v\n", err)
			return nil, err
		}

		if err := p.DB.UpdateUser(*u); err != nil {
			return nil, err
		}
	}

	return spotify.GetNowPlaying(u.Tokens)
}

func formatSSE(event string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}

	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)

	if err := enc.Encode(m); err != nil {
		return "", err
	}

	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("event: %s\n", event))
	b.WriteString(fmt.Sprintf("data: %v\n\n", buf.String()))

	return b.String(), nil
}
