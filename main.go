package main

import (
	// "log"

	"net/http"
	"time"
	"warehouse/pkg/api"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
)

func main() {
	// implement 3 endpoints using net.http
	r := api.NewRouter()
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// log.Println("db,err", db, err)

	const (
		serverAddr          = ":8080"
		shutdownGracePeriod = 30 * time.Second
	)
	srv := &http.Server{
		Addr:              serverAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	//graceful shut down
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	waitForGracefulShutdown(errCh, srv, time.Second*30)

}
