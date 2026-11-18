package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/auth"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/chat"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/database"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/llm"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/recurring"
	"github.com/rafabcanedo/chronos-agent/chronos-golang/internal/transactions"
	"github.com/rafabcanedo/obsidian-golang-toolkit/api/jwtauth"
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
	abacusKey   := env.GetString("ABACUS_API_KEY", "")
	abacusModel := env.GetString("ABACUS_MODEL", "gpt-4o-mini")

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

	recurringSvc    := recurring.NewService(db)
	recurringHandler := recurring.NewHandler(recurringSvc)

	llmClient  := llm.NewClient(abacusKey, abacusModel)
	chatSvc    := chat.NewService(txSvc, recurringSvc, llmClient)
	chatHandler := chat.NewHandler(chatSvc)

	protected := jwtauth.Middleware(jwtSecret)

	mux := http.NewServeMux()

	mux.Handle("GET /health", health.New().Handler())
	mux.HandleFunc("POST /api/auth/user",    authHandler.Register)
	mux.HandleFunc("POST /api/auth/login",   authHandler.Login)
	mux.Handle("PATCH /api/auth/user",       protected(http.HandlerFunc(authHandler.Update)))
	mux.Handle("DELETE /api/auth/user",      protected(http.HandlerFunc(authHandler.Delete)))

	mux.Handle("GET /api/transactions/summary",  protected(http.HandlerFunc(txHandler.Summary)))
	mux.Handle("GET /api/transactions",          protected(http.HandlerFunc(txHandler.List)))
	mux.Handle("POST /api/transactions",         protected(http.HandlerFunc(txHandler.Create)))
	mux.Handle("PATCH /api/transactions/{id}",   protected(http.HandlerFunc(txHandler.Update)))
	mux.Handle("DELETE /api/transactions/{id}",  protected(http.HandlerFunc(txHandler.Delete)))

	mux.Handle("GET /api/recurring",           protected(http.HandlerFunc(recurringHandler.List)))
	mux.Handle("POST /api/recurring",          protected(http.HandlerFunc(recurringHandler.Create)))
	mux.Handle("PUT /api/recurring/{id}",      protected(http.HandlerFunc(recurringHandler.Update)))
	mux.Handle("DELETE /api/recurring/{id}",   protected(http.HandlerFunc(recurringHandler.Delete)))

	mux.Handle("POST /api/chat", protected(http.HandlerFunc(chatHandler.Handle)))

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
