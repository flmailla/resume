package handlers

import (
	"io"
	"log/slog"
	"testing"

	"github.com/flmailla/resume/logger"
)

func TestMain(m *testing.M) {
	logger.Logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	m.Run()
}
