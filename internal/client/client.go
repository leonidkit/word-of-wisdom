package client

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/leonidkit/word-of-wisdom/internal/messages"
	messagesmodels "github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

var (
	componentName = "tcp-client"

	readTimeoutDefault       = 1 * time.Second
	writeTimeoutDefault      = 1 * time.Second
	keepAliveDurationDefault = 1 * time.Second
)

type messageHandler interface {
	Handle(ctx context.Context, conn net.Conn, msg messagesmodels.Message) error
}

type options struct {
	serverAddr string
	handler    messageHandler

	wg *sync.WaitGroup

	lg *slog.Logger

	readTimeout       time.Duration
	writeTimeout      time.Duration
	keepAliveDuration time.Duration
}

type optionFunc func(opts *options)

func WithReadTimeout(t time.Duration) optionFunc {
	return func(opts *options) {
		opts.readTimeout = t
	}
}

func WithWriteTimeout(ka time.Duration) optionFunc {
	return func(opts *options) {
		opts.keepAliveDuration = ka
	}
}

func WithKeepAlive(t time.Duration) optionFunc {
	return func(opts *options) {
		opts.writeTimeout = t
	}
}

type Client struct {
	options
}

func NewClient(
	serverAddr string,
	handler messageHandler,
	logger *slog.Logger,
	opts ...optionFunc,
) *Client {
	o := options{
		handler:           handler,
		serverAddr:        serverAddr,
		wg:                new(sync.WaitGroup),
		lg:                logger.WithGroup(componentName),
		readTimeout:       readTimeoutDefault,
		writeTimeout:      writeTimeoutDefault,
		keepAliveDuration: keepAliveDurationDefault,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &Client{options: o}
}

func (s *Client) Run(ctx context.Context) error {
	serverAddr, err := net.ResolveTCPAddr("tcp", s.serverAddr)
	if err != nil {
		return fmt.Errorf("resolve server addr: %v", err)
	}

	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return fmt.Errorf("%s: listener creating: %v", componentName, err)
	}
	defer func() {
		if err := conn.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			s.lg.Error("close error", slog.Any("error", err))
		}
	}()

	s.lg.Info("dial and send")

	if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
		return fmt.Errorf("set read deadline: %v", err)
	}
	if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
		return fmt.Errorf("set write deadline: %v", err)
	}

	req, err := messages.Encode{}.Encode(messagesmodels.NewQuoteRequestMessage())
	if err != nil {
		return fmt.Errorf("encode request message: %v", err)
	}

	_, err = conn.Write(append(req, '\n'))
	if err != nil {
		return fmt.Errorf("write message: %v", err)
	}

	for i := 0; i < 2; i++ {
		reader := bufio.NewReader(conn)
		data, err := reader.ReadBytes('\n')
		if err != nil {
			return fmt.Errorf("connection read: %v", err)
		}

		d, err := messages.Decoder{}.Decode(bytes.TrimSpace(data))
		if err != nil {
			return fmt.Errorf("message decode: %v", err)
		}

		if err := s.handler.Handle(ctx, conn, d); err != nil {
			return fmt.Errorf("message handle: %v", err)
		}
	}

	return nil
}
