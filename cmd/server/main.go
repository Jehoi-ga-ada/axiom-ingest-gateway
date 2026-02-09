package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := config.NewLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %s\n", err.Error())
		return
	}

	viper, err := config.NewViper()
	if err != nil {
		logger.Fatal("Viper failed to load",
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
		MaxHeaderBytes: 1024 * 1024,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		logger.Info("Shutdown signal received")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("HTTP server Shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("HTTP server ListenAndServe", zap.Error(err))
	}

	<-idleConnsClosed
}
