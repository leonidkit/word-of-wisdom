package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/leonidkit/word-of-wisdom/internal/client"
	"github.com/leonidkit/word-of-wisdom/internal/logger"
	"github.com/leonidkit/word-of-wisdom/internal/messages"
	"github.com/leonidkit/word-of-wisdom/internal/messageshandler"
	messagesmodel "github.com/leonidkit/word-of-wisdom/internal/models/messages"
	challenger "github.com/leonidkit/word-of-wisdom/internal/services/client-challenger"
)

var appName = "word-of-wisdom"

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Logger.
	logger.MustInit(logger.Options{
		Level:          "debug",
		EnableSource:   false,
		ProductionMode: true,
	})

	// Services.
	serverChallengerSvc := challenger.New()

	// Handlers.
	handlers := client.NewHandlers(
		appName,
		serverChallengerSvc,
		messages.Decoder{},
		messages.Encode{},
	)

	mh := messageshandler.New()
	err := mh.RegisterHandler(
		new(messagesmodel.ChallengeRequestMessage).Name(),
		messageshandler.HandlerFunc(handlers.HandleChallengeRequestMessage),
	)
	if err != nil {
		return fmt.Errorf("register `HandleChallengeRequestMessage` handle: %v", err)
	}

	err = mh.RegisterHandler(
		new(messagesmodel.QuoteResponseMessage).Name(),
		messageshandler.HandlerFunc(handlers.HandleQuoteResponseMessage),
	)
	if err != nil {
		return fmt.Errorf("register `HandleQuoteResponseMessage` handle: %v", err)
	}

	// Server.
	srv := client.NewClient(":8080", mh, slog.Default())

	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
