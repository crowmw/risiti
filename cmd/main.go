package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crowmw/risiti/internal/filestore"
	"github.com/crowmw/risiti/internal/handlers"
	m "github.com/crowmw/risiti/internal/middleware"
	"github.com/crowmw/risiti/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	PORT = ":2137"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	filestore := filestore.NewFileStore()
	store := store.NewStore()
	fileserver := http.FileServer(http.Dir("static"))
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.Recoverer, m.CSPMiddleware)

	router.Handle("/static/*", http.StripPrefix("/static/", fileserver))

	router.Get("/", handlers.NewGetHomeHandler().ServeHTTP)
	router.Get("/upload", handlers.NewGetUploadHandler().ServeHTTP)
	router.Post("/submit", handlers.NewPostSubmitHandler(filestore, store).ServeHTTP)

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", PORT))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
