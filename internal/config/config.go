package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Mode *string
	Port *int
	Host *string
	Env  *Environment
}

type Environment struct {
	DatabaseURL         string
	SpotifyClientID     string
	SpotifyClientSecret string
}

func InitConfig() *Config {
	mode := flag.String("mode", "production", "Sets the mode of program execution.")
	port := flag.Int("port", 3000, "Sets the port for the server to listen on.")
	host := flag.String("host", "localhost", "Sets the host for the server to listen on.")
	flag.Parse()

	env := getEnv()

	return &Config{
		Mode: mode,
		Port: port,
		Host: host,
		Env:  env,
	}
}

func getEnv() *Environment {
	var env = &Environment{}
	required := []string{"DATABASE_URL", "SPOTIFY_CLIENT_ID", "SPOTIFY_CLIENT_SECRET"}

	for _, key := range required {
		v := os.Getenv(key)
		if v == "" {
			log.Fatalf("Environment variable %s is required.", key)
		} else {
			switch key {
			case "DATABASE_URL":
				env.DatabaseURL = v
			case "SPOTIFY_CLIENT_ID":
				env.SpotifyClientID = v
			case "SPOTIFY_CLIENT_SECRET":
				env.SpotifyClientSecret = v
			}
		}
	}

	return env
}
