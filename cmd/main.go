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
	receiptrepo "github.com/crowmw/risiti/internal/repo"
	"github.com/crowmw/risiti/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	PORT = ":2137"
)

func main() {
	filestore := filestore.NewFileStore()
	store := store.NewStore()
	receiptRepo := receiptrepo.NewReceiptRepo(receiptrepo.Receipt{}, store)
	fileserver := http.FileServer(http.Dir("static"))
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.Recoverer, m.CSPMiddleware)

	router.Handle("/static/*", http.StripPrefix("/static/", fileserver))

	router.Get("/", handlers.NewGetHomeHandler().ServeHTTP)
	router.Get("/receipts", handlers.NewGetReceiptsHandler(receiptRepo).ServeHTTP)
	router.Get("/upload", handlers.NewGetUploadHandler().ServeHTTP)
	router.Post("/submit", handlers.NewPostSubmitHandler(filestore, receiptRepo).ServeHTTP)
	router.Post("/search", handlers.NewPostSearchHandler(receiptRepo).ServeHTTP)

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

	slog.Info("🚀 Server started! Listening on port " + PORT)
	<-killSig

	slog.Info("🚨 Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}
