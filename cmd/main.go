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

	"github.com/crowmw/risiti/internal/handlers"
	m "github.com/crowmw/risiti/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()
	r.Use(middleware.Logger, m.CSPMiddleware)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", handlers.HomeGetHandler)
	r.Get("/upload", handlers.UploadHandler)
	r.Post("/submit", handlers.SubmitHandler)

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	port := ":2137"
	srv := &http.Server{
		Addr:    port,
		Handler: r,
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

	logger.Info("Server started", slog.String("port", port))
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
