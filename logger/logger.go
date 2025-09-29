package logger

import (
	"io"
	"log/slog"
	"os"
	"gopkg.in/natefinch/lumberjack.v2"
)

var     Logger *slog.Logger
const   loggerPath string = "/var/log/resume/app.log"

func InitLogger() error {

	logRotator := &lumberjack.Logger{
		Filename:   loggerPath,
		MaxSize:    100,   
		MaxBackups: 5,     
		MaxAge:     30,    
		Compress:   true, 
	}

	multiWriter := slog.NewJSONHandler(
		io.MultiWriter(os.Stdout, logRotator),
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		},
	)

	Logger = slog.New(multiWriter)
	
	slog.SetDefault(Logger)

	return nil
}