package main

import (
	"context"
	"log"

	"yoyo-server/internal/config"
	"yoyo-server/internal/database"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/seeds"
	"yoyo-server/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	repos := repositories.New(db)
	svc := services.New(cfg, db, repos)
	if err := seeds.Run(context.Background(), cfg, db, svc); err != nil {
		log.Fatal(err)
	}
	log.Println("seed completed")
}
