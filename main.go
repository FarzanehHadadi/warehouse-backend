package main

import (
	"log"
	"net/http"
	"time"

	"warehouse/pkg/api"
	"warehouse/pkg/database"
	"warehouse/pkg/repository"
)

//	@title			Warehouse API
//	@version		1.2
//	@description	Warehouse Management System API

//	@host		localhost:8080
//	@BasePath	/v1

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func main() {
	// Database
	dbConfig := database.DefaultConfig()
	db, err := database.NewPostgresConfiguration(&dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Repository (now cleanly initialized in its own package)
	repo := repository.NewRepository(db)

	// API Router with handler
	router := api.NewRouter(repo)

	// HTTP Server
	const (
		serverAddr          = ":8080"
		shutdownGracePeriod = 30 * time.Second
	)

	srv := &http.Server{
		Addr:              serverAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("🚀 Server starting on %s", serverAddr)
		errCh <- srv.ListenAndServe()
	}()

	waitForGracefulShutdown(errCh, srv, shutdownGracePeriod)
}

// Graceful shutdown function (kept the same)
