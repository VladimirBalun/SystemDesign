package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	storageAddress := os.Getenv("STORAGE_ADDRESS")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storage, err := NewStorage(ctx, storageAddress)
	if err != nil {
		log.Fatalf("failed to connect to storage: %s", err.Error())
		return
	}

	router := mux.NewRouter()
	handler := NewHandler(&storage)
	router.HandleFunc("/lock/{resource}", handler.Lock).Methods(http.MethodPost)
	router.HandleFunc("/unlock/{resource}", handler.Unlock).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	<-ctx.Done()

	const shutdownTimeout = 3 * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("server gracefull shutdown failed: %s", err.Error())
	}
}
