package messageshandler

import (
	"context"
	"errors"
	"net"

	"github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

var (
	ErrHandlerExists   = errors.New("handler already registered")
	ErrHandlerNotFound = errors.New("handler not found")
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn, msg messages.Message) error
}

type HandlerFunc func(ctx context.Context, conn net.Conn, msg messages.Message) error

func (hf HandlerFunc) Handle(ctx context.Context, conn net.Conn, msg messages.Message) error {
	return hf(ctx, conn, msg)
}

type MessageHandler struct {
	handlers map[string]Handler
}

func New() MessageHandler {
	return MessageHandler{
		handlers: make(map[string]Handler),
	}
}

// RegisterHandler registers handler by event name.
func (mh MessageHandler) RegisterHandler(messageName string, h Handler) error {
	if _, ok := mh.handlers[messageName]; ok {
		return ErrHandlerExists
	}
	mh.handlers[messageName] = h

	return nil
}

func (mh MessageHandler) Handle(ctx context.Context, conn net.Conn, msg messages.Message) error {
	if h, ok := mh.handlers[msg.Name()]; ok {
		return h.Handle(ctx, conn, msg)
	}
	return ErrHandlerNotFound
}
