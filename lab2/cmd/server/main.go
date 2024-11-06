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

	"communicator/handlers/fileupload"
	"communicator/handlers/moviecrud"
)

func main() {
	// conf, err := config.FromEnv()
	// if err != nil {
	// 	log.Fatal("Load config from environment:", err)
	// }

	movieCRUDGroup := moviecrud.HandlerGroup{}
	fileUploadGroup := fileupload.HandlerGroup{}

	serverMux := http.NewServeMux()
	serverMux.Handle("/crud", movieCRUDGroup.Mux())
	serverMux.Handle("/file", fileUploadGroup.Mux())

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: serverMux,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		<-sig
		return errors.New("received sigterm")
	})

	eg.Go(func() error {
		err := server.ListenAndServe()
		return fmt.Errorf("serve moviecrud: %w", err)
	})

	eg.Go(func() error {
		<-ctx.Done()
		log.Println("Stopping")

		cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err := server.Shutdown(cleanupCtx)
		if err != nil {
			log.Print("Error closing server: ", err)
		}
		cancel()
		return nil
	})

	err := eg.Wait()
	log.Println("Stop reason:", err)
}
