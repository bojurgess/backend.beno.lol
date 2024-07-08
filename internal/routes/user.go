package routes

import (
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
	DB       *database.Database
	Config   *config.Config
	Channels map[string]*Channel
}

type Channel struct {
	Subscribers int
	NowPlaying  chan spotify.NowPlayingResponse
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

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// if channel doesn't already exist we make it and start polling from Spotify every second
	ch := p.Channels[u.ID]

	if ch == nil {
		ch = &Channel{
			Subscribers: 1,
			NowPlaying:  make(chan spotify.NowPlayingResponse),
		}
		p.Channels[u.ID] = ch

		// loop every second
		go func() {
			for {
				if ch.Subscribers == 0 {
					delete(p.Channels, u.ID)
					return
				}
				p.NowPlaying(&u, nil)
				time.Sleep(1 * time.Second)
			}
		}()
	} else {
		ch.Subscribers++
	}

	// send data every second, until client disconnects
	for {
		select {
		case <-ctx.Done():
			ch.Subscribers--
			fmt.Println("Client disconnected")
			return
		case <-time.After(1 * time.Second):
			// np, err := p.getNowPlaying(&u)
			// if err != nil {
			// 	fmt.Println(err)
			// 	continue
			// }
			np := <-ch.NowPlaying

			b, err := json.Marshal(np)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Fprintf(w, "data: %v\n\n", string(b))

			w.(http.Flusher).Flush()
		}
	}
}

func (p *User) NowPlaying(u *database.User, last *spotify.NowPlayingResponse) {
	fmt.Printf("Polling now playing for %s\n", u.ID)
	ch := p.Channels[u.ID]

	if time.Now().After(u.Tokens.ExpiresAt) {
		if err := spotify.RefreshAccessToken(&u.Tokens, p.Config.Env.SpotifyClientID, p.Config.Env.SpotifyClientSecret); err != nil {
			fmt.Printf("Error refreshing token: %v\n", err)
			return
		}

		if err := p.DB.UpdateUser(*u); err != nil {
			fmt.Printf("Error updating user: %v\n", err)
			return
		}
	}

	np, err := spotify.GetNowPlaying(u.Tokens)
	if err != nil {
		fmt.Printf("Error getting now playing: %v\n", err)
		return
	}

	ch.NowPlaying <- *np
}
