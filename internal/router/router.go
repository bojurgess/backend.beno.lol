package router

import (
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
)

type Router struct {
	DB     *database.Database
	Config *config.Config
	Mux    http.Handler
}

func Create(app config.Application) *Router {
	r := &Router{
		DB:     app.DB,
		Config: app.Config,
		Mux:    http.NewServeMux(),
	}

	return r
}
