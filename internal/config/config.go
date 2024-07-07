package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Mode *string
	Env  *Environment
}

type Environment struct {
	DatabaseURL         string
	SpotifyClientID     string
	SpotifyClientSecret string
}

func InitConfig() *Config {
	mode := flag.String("mode", "production", "Sets the mode of program execution.")
	flag.Parse()

	env := getEnv()

	return &Config{
		Mode: mode,
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
