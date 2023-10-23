package server

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

	"golang.org/x/sync/errgroup"

	"github.com/leonidkit/word-of-wisdom/internal/messages"
	messagesmodels "github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

var (
	componentName = "tcp-server"

	readTimeoutDefault       = 1 * time.Second
	writeTimeoutDefault      = 1 * time.Second
	keepAliveDurationDefault = 1 * time.Second
)

type messageHandler interface {
	Handle(ctx context.Context, conn net.Conn, msg messagesmodels.Message) error
}

type options struct {
	addr    string
	handler messageHandler

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

type Server struct {
	options
}

func NewServer(
	addr string,
	handler messageHandler,
	logger *slog.Logger,
	opts ...optionFunc,
) *Server {
	o := options{
		handler:           handler,
		addr:              addr,
		wg:                new(sync.WaitGroup),
		lg:                logger.WithGroup(componentName),
		readTimeout:       readTimeoutDefault,
		writeTimeout:      writeTimeoutDefault,
		keepAliveDuration: keepAliveDurationDefault,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &Server{options: o}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("%s: listener creating: %v", componentName, err)
	}
	defer func() {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			s.lg.Error("close error", slog.Any("error", err))
		}
	}()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		if err := listener.Close(); err != nil {
			return fmt.Errorf("listener close: %v", err)
		}

		s.wg.Wait()
		return nil
	})

	eg.Go(func() error {
		s.lg.Info("listen and accept", slog.String("addr", s.addr))

		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return nil
				}
				return fmt.Errorf("%s: accept connection: %v", componentName, err)
			}

			if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
				return fmt.Errorf("%s: set read deadline: %v", componentName, err)
			}
			if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
				return fmt.Errorf("%s: set write deadline: %v", componentName, err)
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				if err := s.processConnection(ctx, conn); err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					s.lg.Error("process connection", slog.Any("error", err))
				}
			}()
		}
	})

	return eg.Wait()
}

func (s *Server) processConnection(ctx context.Context, conn net.Conn) error {
	s.lg.Debug("connection remote addr", slog.String("addr", conn.RemoteAddr().String()))

	defer func() {
		if err := conn.Close(); err != nil {
			s.lg.Error("connection close", slog.Any("error", err))
		}
	}()

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
	return nil
}
