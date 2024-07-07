package main

import (
	"fmt"
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/router"
)

func main() {
	app := config.Application{
		DB:     &database.Database{},
		Config: config.InitConfig(),
	}

	app.DB.Connect(app.Config.Env.DatabaseURL)
	r := router.Create(app)

	addr := fmt.Sprintf("http://%s:%d", *app.Config.Host, *app.Config.Port)
	fmt.Printf("Server listening on %s\n", addr)

	http.ListenAndServe(addr, r.Mux)
}
