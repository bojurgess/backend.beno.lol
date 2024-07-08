package main

import (
	"fmt"
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
	"github.com/bojurgess/backend.beno.lol/internal/router"
	"github.com/rs/cors"
)

func main() {
	app := config.Application{
		DB:     &database.Database{},
		Config: config.InitConfig(),
	}

	app.DB.Connect(app.Config.Env.DatabaseURL)
	r := router.Create(app)

	handler := cors.Default().Handler(r.Mux)

	addr := fmt.Sprintf(":%d", *app.Config.Port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		panic(err)
	}
}
