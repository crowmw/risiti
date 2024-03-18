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

	"github.com/crowmw/risiti/handler"
	m "github.com/crowmw/risiti/middleware"
	"github.com/crowmw/risiti/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	PORT = ":2137"
)

func main() {
	fileserver := http.FileServer(http.Dir("static"))

	// Services
	fs := service.NewFileStorage()
	db := service.NewDB()
	receiptService := service.NewReceiptService(db)

	// Handlers
	basicHandler := handler.NewBasicHandler(receiptService)
	receiptHandler := handler.NewReceiptHandler(receiptService, fs)

	// Routes
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer, m.CSPMiddleware)

	router.Handle("/static/*", http.StripPrefix("/static/", fileserver))

	// Views
	router.Get("/", basicHandler.GetHome)
	router.Get("/upload", basicHandler.GetUpload)

	// Partials
	router.Get("/receipts", receiptHandler.GetReceipts)
	router.Post("/receipt", receiptHandler.PostReceipt)
	router.Post("/search", receiptHandler.SearchReceipts)

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
