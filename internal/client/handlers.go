package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	messagesmodels "github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

var ErrWrongMessageType = errors.New("wrong message type")

type challengerService interface {
	Accept(ctx context.Context, challenge string) (string, error)
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
	messageDecoder messageDecoder
	messageEncoder messageEncoder
}

func NewHandlers(
	domain string,
	challengerSvc challengerService,
	messageDecoder messageDecoder,
	messageEncoder messageEncoder,
) *Handlers {
	return &Handlers{
		domain:         domain,
		challengerSvc:  challengerSvc,
		messageDecoder: messageDecoder,
		messageEncoder: messageEncoder,
	}
}

func (h *Handlers) HandleChallengeRequestMessage(ctx context.Context, conn net.Conn, msg messagesmodels.Message) error {
	v, ok := msg.(*messagesmodels.ChallengeRequestMessage)
	if !ok {
		return ErrWrongMessageType
	}

	solution, err := h.challengerSvc.Accept(ctx, v.Challenge)
	if err != nil {
		return fmt.Errorf("`handle challenge request` challenge accept: %v", err)
	}

	data, err := h.messageEncoder.Encode(messagesmodels.NewChallengeResponseMessage(solution))
	if err != nil {
		return fmt.Errorf("`handle challenge request` encode solution: %v", err)
	}

	err = h.writeMsg(conn, data)
	if err != nil {
		return fmt.Errorf("write solution: %v", err)
	}

	return nil
}

func (h *Handlers) HandleQuoteResponseMessage(ctx context.Context, conn net.Conn, msg messagesmodels.Message) error {
	v, ok := msg.(*messagesmodels.QuoteResponseMessage)
	if !ok {
		return ErrWrongMessageType
	}

	slog.Info(v.Quote)

	return nil
}

func (h *Handlers) writeMsg(conn net.Conn, msg []byte) (err error) {
	_, err = conn.Write(append(msg, '\n'))
	return
}
