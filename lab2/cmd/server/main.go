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
	"communicator/handlers/ws"
)

func main() {
	// conf, err := config.FromEnv()
	// if err != nil {
	// 	log.Fatal("Load config from environment:", err)
	// }

	movieCRUDGroup := moviecrud.HandlerGroup{}
	fileUploadGroup := fileupload.HandlerGroup{}

	httpMux := http.NewServeMux()
	httpMux.Handle("/crud", movieCRUDGroup.Mux())
	httpMux.Handle("/file", fileUploadGroup.Mux())

	httpServer := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: httpMux,
	}

	wsGroup := ws.NewHandlerGroup()
	wsMux := wsGroup.Mux()

	wsServer := http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: wsMux,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		<-sig
		return errors.New("received sigterm")
	})

	eg.Go(func() error {
		err := httpServer.ListenAndServe()
		return fmt.Errorf("serve http: %w", err)
	})

	eg.Go(func() error {
		err := wsServer.ListenAndServe()
		return fmt.Errorf("serve ws: %w", err)
	})

	eg.Go(func() error {
		err := wsGroup.HandleMessages(ctx)
		return fmt.Errorf("handle ws messages: %w", err)
	})

	eg.Go(func() error {
		<-ctx.Done()
		log.Println("Stopping")

		cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err := httpServer.Shutdown(cleanupCtx)
		if err != nil {
			log.Print("Error closing http server: ", err)
		}

		err = wsServer.Shutdown(cleanupCtx)
		if err != nil {
			log.Print("Error closing ws server: ", err)
		}

		cancel()
		return nil
	})

	err := eg.Wait()
	log.Println("Stop reason:", err)
}
