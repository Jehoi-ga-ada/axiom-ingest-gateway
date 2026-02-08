package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := config.NewLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %s\n", err.Error())
	}

	viper, err := config.NewViper()
	if err != nil {
		logger.Fatal("Viper failed to unload",
			zap.String("err", err.Error()),
		)
	}

	router := config.NewRouter(logger)

	config.NewApp(&config.Config{
		Router: router,
		Viper: viper,
		Logger: logger,
	})

	port := viper.GetString("PORT")
	if port == "" {
		logger.Info("PORT not found, defaulting to :8000")
		port = "8000"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	logger.Info("Starting server on",
		zap.String("port", port),
	)

	server := &http.Server{
		Addr: serverAddr,
		Handler: router,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal("Server failed to start",
			zap.String("err", err.Error()),
		)
	}
}
