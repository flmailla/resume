package handlers

import (
	"testing"
	"log/slog"
	"github.com/flmailla/resume/logger"
	"io"
)

func TestMain(m *testing.M) {
	logger.Logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	m.Run()
}
