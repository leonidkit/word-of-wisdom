package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v9"

	"github.com/leonidkit/word-of-wisdom/internal/config"
	"github.com/leonidkit/word-of-wisdom/internal/logger"
	"github.com/leonidkit/word-of-wisdom/internal/messages"
	"github.com/leonidkit/word-of-wisdom/internal/messageshandler"
	messagesmodel "github.com/leonidkit/word-of-wisdom/internal/models/messages"
	"github.com/leonidkit/word-of-wisdom/internal/repositories/challenges"
	"github.com/leonidkit/word-of-wisdom/internal/repositories/quotes"
	"github.com/leonidkit/word-of-wisdom/internal/server"
	challenger "github.com/leonidkit/word-of-wisdom/internal/services/server-challenger"
)

var appName = "word-of-wisdom"

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() error {
	cfg := config.ServerConfig{}

	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("config parse: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Logger.
	logger.MustInit(logger.Options{
		Level:          cfg.LogLevel,
		ProductionMode: cfg.IsProduction(),
	})

	// Repositories.
	challengesRepo := challenges.New()
	quotesRepo := quotes.New()

	// Services.
	serverChallengerSvc := challenger.New(challengesRepo)

	// Handlers.
	handlers := server.NewHandlers(
		appName,
		serverChallengerSvc,
		messages.Decoder{},
		messages.Encode{},
		quotesRepo,
	)

	mh := messageshandler.New()
	err = mh.RegisterHandler(
		new(messagesmodel.ChallengeResponseMessage).Name(),
		messageshandler.HandlerFunc(handlers.HandleChallengeResponseMessage),
	)
	if err != nil {
		return fmt.Errorf("register `HandleChallengeResponseMessage` handle: %v", err)
	}

	err = mh.RegisterHandler(
		new(messagesmodel.QuoteRequestMessage).Name(),
		messageshandler.HandlerFunc(handlers.HandleQuoteRequestMessage),
	)
	if err != nil {
		return fmt.Errorf("register `HandleQuoteRequestMessage` handle: %v", err)
	}

	// Server.
	srv := server.NewServer(cfg.Addr, mh, slog.Default())

	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
