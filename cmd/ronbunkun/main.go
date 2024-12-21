package main

import (
	"log"
	config "ronbunkun/config"
	"ronbunkun/server"
)

func main() {
	cfg := config.NewConfig()

	app := server.NewServer(cfg)

	server.ConfigureRoutes(app)

	err := app.Start(cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Port already used")
	}
}
