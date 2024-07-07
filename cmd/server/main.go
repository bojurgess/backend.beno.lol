package main

import (
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
}
