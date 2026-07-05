package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"warehouse/pkg/api"
	"warehouse/pkg/api/auth"
	"warehouse/pkg/database"
	"warehouse/pkg/events"
	"warehouse/pkg/listeners"
	"warehouse/pkg/logger"
	"warehouse/pkg/repository"

	"github.com/joho/godotenv"
)

//	@title			Warehouse API
//	@version		1.2
//	@description	Warehouse Management System API

//	@host		localhost:8080
//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
//	@description				Enter your API Key here

//	@securityDefinitions.apikey	AdminRegistrationKeyAuth
//	@in							header
//	@name						X-Admin-Key
//	@description				Enter your Admin Registration Key here

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Enter JWT with Bearer prefix
func main() {
	logger.Init()
	defer logger.Sync()
	// env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database
	dbConfig := database.DefaultConfig()
	db, err := database.NewPostgresConfiguration(&dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Repository (now cleanly initialized in its own package)
	repo := repository.NewRepository(db)
	setupEventListeners(repo)
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	auth.LoadSecrets([]byte(jwtSecretKey), []byte(refreshSecretKey))
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

// Graceful shutdown func
// tion (kept the same)
func setupEventListeners(repo *repository.Repository) {
	logger := listeners.NewActivityLogger(repo.Activity) // using your ActivityRepository

	events.Bus.Subscribe(events.Created, logger.Handle)
	events.Bus.Subscribe(events.Updated, logger.Handle)
	events.Bus.Subscribe(events.Deleted, logger.Handle)
}
