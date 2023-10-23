package server

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

	messagesmodels "github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

var ErrWrongMessageType = errors.New("wrong message type")

type ctxValidateKey struct{}

type challengerService interface {
	NewChallenge(ctx context.Context, domain string) (string, error)
	Validate(ctx context.Context, challenge string) (bool, error)
}

type quotesRepository interface {
	GetRandom(ctx context.Context) (string, error)
}

type messageDecoder interface {
	Decode(msg []byte) (messagesmodels.Message, error)
}

type messageEncoder interface {
	Encode(msg messagesmodels.Message) ([]byte, error)
}

type Handlers struct {
	domain         string
	challengerSvc  challengerService
	quotesRepo     quotesRepository
	messageDecoder messageDecoder
	messageEncoder messageEncoder
}

func NewHandlers(
	domain string,
	challengerSvc challengerService,
	messageDecoder messageDecoder,
	messageEncoder messageEncoder,
	quotesRepo quotesRepository,
) *Handlers {
	return &Handlers{
		domain:         domain,
		challengerSvc:  challengerSvc,
		quotesRepo:     quotesRepo,
		messageDecoder: messageDecoder,
		messageEncoder: messageEncoder,
	}
}

func (h *Handlers) HandleChallengeResponseMessage(ctx context.Context, conn net.Conn, msg messagesmodels.Message) error {
	if err := h.verifyMsg(ctx, msg); err != nil {
		return fmt.Errorf("`handle challenge response message` verify message: %v", err)
	}
	return h.HandleQuoteRequestMessage(context.WithValue(ctx, ctxValidateKey{}, true), conn, msg)
}

func (h *Handlers) HandleQuoteRequestMessage(ctx context.Context, conn net.Conn, msg messagesmodels.Message) error {
	handlerName := "handle quote request"

	if v, ok := ctx.Value(ctxValidateKey{}).(bool); !ok || !v {
		if err := h.challenge(ctx, conn); err != nil {
			return fmt.Errorf("`%s`: challenge: %v", handlerName, err)
		}
	}

	quote, err := h.quotesRepo.GetRandom(ctx)
	if err != nil {
		return fmt.Errorf("`%s` get random quote: %v", handlerName, err)
	}

	data, err := h.messageEncoder.Encode(messagesmodels.NewQuoteResponseMessage(quote))
	if err != nil {
		return fmt.Errorf("`%s` encode quote: %v", handlerName, err)
	}

	err = h.writeMsg(conn, data)
	if err != nil {
		return fmt.Errorf("`%s` write challenge: %v", handlerName, err)
	}
	return nil
}

func (h *Handlers) challenge(ctx context.Context, conn net.Conn) error {
	challenge, err := h.challengerSvc.NewChallenge(ctx, h.domain)
	if err != nil {
		return fmt.Errorf("new challenge: %v", err)
	}

	data, err := h.messageEncoder.Encode(messagesmodels.NewChallengeRequestMessage(challenge))
	if err != nil {
		return fmt.Errorf("encode challenge: %v", err)
	}

	err = h.writeMsg(conn, data)
	if err != nil {
		return fmt.Errorf("write challenge: %v", err)
	}

	reader := bufio.NewReader(conn)
	data, err = reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("connection read: %v", err)
	}

	m, err := h.messageDecoder.Decode(bytes.TrimSpace(data))
	if err != nil {
		return fmt.Errorf("message decode: %v", err)
	}

	if err := h.verifyMsg(ctx, m); err != nil {
		return fmt.Errorf("verify message: %v", err)
	}
	return nil
}

func (h *Handlers) verifyMsg(ctx context.Context, msg messagesmodels.Message) error {
	m, ok := msg.(*messagesmodels.ChallengeResponseMessage)
	if !ok {
		return ErrWrongMessageType
	}

	ok, err := h.challengerSvc.Validate(ctx, m.Solution)
	if err != nil {
		return fmt.Errorf("challenge validate: %v", err)
	}
	if !ok {
		return fmt.Errorf("challenge wrong solution: %v", err)
	}

	return nil
}

func (h *Handlers) writeMsg(conn net.Conn, msg []byte) (err error) {
	_, err = conn.Write(append(msg, '\n'))
	return
}
