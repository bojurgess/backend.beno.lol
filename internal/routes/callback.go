package routes

import (
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
)

type Callback struct {
	DB     *database.Database
	Config *config.Config
}

func (p *Callback) Route(w http.ResponseWriter, r *http.Request) {
	// Implement the callback route here.
	w.Write([]byte("Callback route"))
}
