package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"communicator/config"
	"communicator/connections/db"
	"communicator/handlers/fileupload"
	"communicator/handlers/moviecrud"
	"communicator/handlers/tcp"
	"communicator/handlers/ws"
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

	movieCRUDGroup := moviecrud.HandlerGroup{
		Database: dbConn,
	}
	fileUploadGroup := fileupload.HandlerGroup{
		Directory: conf.Upload.Directory,
	}

	// Set up routing
	httpMux := http.NewServeMux()
	httpMux.Handle("/movie/", http.StripPrefix("/movie", movieCRUDGroup.Mux()))
	httpMux.Handle("/file/", http.StripPrefix("/file", fileUploadGroup.Mux()))

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

	tcpServer := tcp.Server{
		FilePath: conf.TCP.FilePath,
	}
	tcpListener, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		log.Fatalf("licten on 8082: %s", err)
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

	eg.Go(func() error {
		err := wsServer.ListenAndServe()
		return fmt.Errorf("serve ws: %w", err)
	})

	eg.Go(func() error {
		err := wsGroup.HandleMessages(ctx)
		return fmt.Errorf("handle ws messages: %w", err)
	})

	eg.Go(func() error {
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				return fmt.Errorf("accept tcp connection: %w", err)
			}
			go tcpServer.HandleRequest(conn)
		}
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

		err = wsServer.Shutdown(cleanupCtx)
		if err != nil {
			log.Print("Error closing ws server: ", err)
		}

		err = tcpListener.Close()
		if err != nil {
			log.Print("Error closing tcp server: ", err)
		}
		cancel()
		return nil
	})

	err = eg.Wait()
	log.Println("Stop reason:", err)
}
