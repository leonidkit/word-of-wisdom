package logger_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leonidkit/word-of-wisdom/internal/logger"
)

func TestInit(t *testing.T) {
	err := logger.Init(logger.Options{
		Level:          "error",
		EnableSource:   false,
		ProductionMode: true,
	})
	require.NoError(t, err)

	slog.Default().WithGroup("user-cache").Error("inconsistent state", slog.String("uid", "1234"))
	// Output: {"T":"2023-10-11 17:15:46","level":"ERROR","msg":"inconsistent state","user-cache":{"uid":"1234"}}
}
