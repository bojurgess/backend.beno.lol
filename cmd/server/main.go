package main

import (
	"github.com/bojurgess/backend.beno.lol/internal/database"
)

type Application struct {
	db *database.Database
}

func main() {
	app := Application{
		db: &database.Database{},
	}

	app.db.Connect("./.db/database.sqlite3")
}
