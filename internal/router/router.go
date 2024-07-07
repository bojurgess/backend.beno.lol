package router

import (
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/routes"
)

type Router struct {
	DB     *database.Database
	Config *config.Config
	Mux    *http.ServeMux
}

func Create(app config.Application) *Router {
	mux := http.NewServeMux()
	auth := &routes.Auth{
		DB:     app.DB,
		Config: app.Config,
	}
	callback := &routes.Callback{
		DB:     app.DB,
		Config: app.Config,
	}
	user := &routes.User{
		DB:     app.DB,
		Config: app.Config,
	}

	mux.HandleFunc("/auth/", auth.Route)
	mux.HandleFunc("/auth/callback/", callback.Route)
	mux.HandleFunc("/user/{id}/", user.Route)

	r := &Router{
		DB:     app.DB,
		Config: app.Config,
		Mux:    mux,
	}

	return r
}
