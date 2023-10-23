package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var Level = new(slog.LevelVar)

var timeKey = "T"

var ErrUnknownLevel = errors.New("unknown level")

type Options struct {
	Level          string
	EnableSource   bool
	ProductionMode bool
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	err := SetLevel(opts.Level)
	if err != nil {
		return err
	}

	config := &slog.HandlerOptions{
		Level: Level,
		AddSource: func() bool {
			return opts.EnableSource
		}(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(timeKey, a.Value.Time().Format(time.DateTime))
			}
			return a
		},
	}

	handler := func() slog.Handler {
		if opts.ProductionMode {
			return slog.NewJSONHandler(os.Stdout, config)
		}
		return slog.NewTextHandler(os.Stdout, config)
	}()

	slog.SetDefault(slog.New(handler))

	return nil
}

func SetLevel(level string) error {
	err := Level.UnmarshalText([]byte(level))
	if err != nil {
		return fmt.Errorf("level unmarshal: %v", err)
	}
	return nil
}
