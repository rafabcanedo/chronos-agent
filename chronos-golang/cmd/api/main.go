package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/auth"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/database"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/transactions"
	toolkitauth "github.com/rafabcanedo/obsidian-golang-toolkit/api/auth"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/middleware"
	"github.com/rafabcanedo/obsidian-golang-toolkit/config"
	"github.com/rafabcanedo/obsidian-golang-toolkit/env"
	"github.com/rafabcanedo/obsidian-golang-toolkit/health"
	"github.com/rafabcanedo/obsidian-golang-toolkit/util/logutil"
)

func main() {
	if vars, err := config.LoadEnv(".env"); err == nil {
		for k, v := range vars {
			os.Setenv(k, v)
		}
	}

	appEnv      := env.GetString("APP_ENV", "development")
	port        := env.GetString("PORT", "8080")
	databaseURL := env.GetString("DATABASE_URL", "")
	jwtSecret   := env.GetString("JWT_SECRET", "")

	logger := logutil.NewLogger(appEnv)

	ctx := context.Background()

	db, err := database.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("database connected")

	authSvc     := auth.NewService(db, jwtSecret)
	authHandler := auth.NewHandler(authSvc)

	txSvc     := transactions.NewService(db)
	txHandler := transactions.NewHandler(txSvc)

	protected := toolkitauth.Middleware(jwtSecret)

	mux := http.NewServeMux()

	mux.Handle("GET /health", health.New().Handler())
	mux.HandleFunc("POST /api/auth/user",  authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.Handle("GET /api/transactions",         protected(http.HandlerFunc(txHandler.List)))
	mux.Handle("PATCH /api/transactions/{id}",  protected(http.HandlerFunc(txHandler.Update)))
	mux.Handle("DELETE /api/transactions/{id}", protected(http.HandlerFunc(txHandler.Delete)))

	cors := middleware.CORS(middleware.CORSOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := middleware.Logger(
		middleware.Recovery(
			middleware.RequestID(
				cors(mux),
			),
		),
	)

	logger.Info("server starting", "port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
