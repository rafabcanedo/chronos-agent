package main

import (
	"log"
	"net/http"
	"os"

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

	appEnv := env.GetString("APP_ENV", "development")
	port := env.GetString("PORT", "8080")

	logger := logutil.NewLogger(appEnv)

	mux := http.NewServeMux()

	mux.Handle("GET /health", health.New().Handler())

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
