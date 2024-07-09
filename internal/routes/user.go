package routes

import (
	"net/http"
	"strings"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
)

type User struct {
	DB     *database.Database
	Config *config.Config
}

// Handles streaming of user NowPlaying data down to client.
// This is done through SSE.
// Route expects to be handled on a route with path parameter {id}.
func (p *User) Route(w http.ResponseWriter, r *http.Request) {
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
}
