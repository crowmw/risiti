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
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	SECRET := os.Getenv("SECRET")

	fileserver := http.FileServer(http.Dir("static"))
	dataImagesServer := http.FileServer(http.Dir("data"))

	// Services
	fs := service.DefaultFileStorage()
	db := service.DefaultDB()
	authService := service.DefaultAuthService([]byte(SECRET))
	receiptService := service.DefaultReceiptService(db)
	userService := service.DefaultUserService(db)

	// Handlers
	basicHandler := handler.NewBasicHandler(receiptService, userService, authService)
	receiptHandler := handler.NewReceiptHandler(receiptService, fs)
	userHandler := handler.NewUserHandler(userService, authService)

	// Routes
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer, m.CORS, m.CSPMiddleware, jwtauth.Verifier(authService.JWTAuth))

	router.Handle("/static/*", http.StripPrefix("/static/", fileserver))

	// Views
	router.Get("/signin", userHandler.GetSignin)
	router.Get("/signup", userHandler.GetSignup)
	router.Get("/", basicHandler.GetHome)

	// Partials
	router.Post("/user", userHandler.PostUser)
	router.Post("/signin", userHandler.PostSignin)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(m.Authenticator(authService.JWTAuth))
		r.Handle("/data/*", http.StripPrefix("/data/", dataImagesServer))
		r.Get("/upload", basicHandler.GetUpload)
		r.Get("/signout", userHandler.GetSignout)

		// Partials
		r.Get("/receipts", receiptHandler.GetReceipts)
		r.Post("/receipt", receiptHandler.PostReceipt)
		r.Post("/search", receiptHandler.SearchReceipts)
	})

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":80",
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

	slog.Info("🚀 Server started! Listening on port 80")
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
