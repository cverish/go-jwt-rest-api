package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/cverish/go-jwt-rest-api/internal/config"
	"github.com/cverish/go-jwt-rest-api/internal/database"
	"github.com/cverish/go-jwt-rest-api/internal/server"
)

var envpath = flag.String("envpath", ".env", "path for desired environment file. default .env")

func main() {
	flag.Parse()

	cfg, err := config.Load(*envpath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if sqlDB, err := db.DB.DB(); err == nil {
		defer sqlDB.Close()
	} else {
		log.Fatal("Failed to close db", err)
	}

	srv := server.NewHttpServer(cfg, db)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start", err)
	}
}
