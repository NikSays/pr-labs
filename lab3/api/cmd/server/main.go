package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"communicator/config"
	"communicator/connections/db"
	"communicator/handlers/monitorcrud"
)

func main() {
	// Load config
	conf, err := config.FromEnv()
	if err != nil {
		log.Fatal("Load config from environment:", err)
	}

	// Prepare dependencies
	dbConn, err := db.NewClient(db.Config{
		Host:     conf.DB.Host,
		Database: conf.DB.Database,
		Username: conf.DB.Username,
		Password: conf.DB.Password,
	})
	if err != nil {
		log.Fatalf("connect to db: %s", err)
	}

	movieCRUDGroup := monitorcrud.HandlerGroup{
		Database: dbConn,
	}

	// Set up routing
	httpMux := http.NewServeMux()
	httpMux.Handle("/monitor/", http.StripPrefix("/monitor", movieCRUDGroup.Mux()))

	httpServer := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: httpMux,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	// Error group. An easy way of creating threads with error handling
	eg, ctx := errgroup.WithContext(context.Background())

	// Handle signals
	eg.Go(func() error {
		<-sig
		return errors.New("received sigterm")
	})

	eg.Go(func() error {
		err := httpServer.ListenAndServe()
		return fmt.Errorf("serve http: %w", err)
	})

	// Cleanup after a thread has returned an error
	eg.Go(func() error {
		<-ctx.Done()
		log.Println("Stopping")

		cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err := httpServer.Shutdown(cleanupCtx)
		if err != nil {
			log.Print("Error closing http server: ", err)
		}

		cancel()
		return nil
	})

	err = eg.Wait()
	log.Println("Stop reason:", err)
}
