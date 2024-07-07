package main

import (
	"fmt"
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
)

type Application struct {
	db     *database.Database
	config *config.Config
}

func main() {
	app := Application{
		db:     &database.Database{},
		config: config.InitConfig(),
	}

	app.db.Connect(app.config.Env.DatabaseURL)

	addr := fmt.Sprintf("http://%s:%d", *app.config.Host, *app.config.Port)
	fmt.Printf("Server listening on %s\n", addr)

	http.ListenAndServe(addr, nil)
}
