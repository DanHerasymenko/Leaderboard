// @title Leaderboard API
// @version 1.0
// @description This is a sample server for Leaderboard API.

package main

import (
	"Leaderboard/cmd/server/handlers"
	"Leaderboard/internal/client"
	"Leaderboard/internal/config"
	"Leaderboard/internal/logger"
	"Leaderboard/internal/services"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err))
	}

	clnts, err := client.NewClients(ctx, cfg)
	if err != nil {
		logger.Panic(ctx, fmt.Errorf("failed to create clients: %w", err))
	}

	svcs := services.NewServices(cfg, clnts)

	app := fiber.New()

	hdlrs := handlers.NewHandlers(cfg, clnts, svcs)
	hdlrs.RegisterRoutes(app)

	go func() {
		err = app.Listen(":8082")
		if err != nil {
			logger.Panic(ctx, fmt.Errorf("failed to start server: %w", err))
		}
	}()

	logger.Info(ctx, "server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	err = app.Shutdown()
}
