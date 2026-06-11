package main

import (
	"context"
	"log"
	"os"

	"yoyo-server/internal/config"
	"yoyo-server/internal/database"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/routes"
	"yoyo-server/internal/seeds"
	"yoyo-server/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	repos := repositories.New(db)
	svc := services.New(cfg, db, repos)

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := seeds.Run(context.Background(), cfg, db, svc); err != nil {
			log.Fatal(err)
		}
		log.Println("seed completed")
		return
	}

	if err := svc.Auth.EnsureSuperAdmin(context.Background()); err != nil {
		log.Fatal(err)
	}
	if _, err := svc.Settings.Get(context.Background()); err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Fatal(err)
	}
	routes.Setup(router, cfg, db, repos, svc)

	log.Printf("YOYO API listening on %s:%s", cfg.Host, cfg.Port)
	if err := router.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
