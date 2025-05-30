package logger

import (
	"log/slog"
	"os"
	"strings"
	"to-do/internal/config"
)

type Logger struct {
	Server *slog.Logger
	Sqlite *slog.Logger
	MW     *slog.Logger
}

func New(cfg *config.Config) (*Logger, error) {
	dir := strings.Split(cfg.Logger.Server, "/")[0]
	os.Mkdir(dir, 0755)

	file1, err := os.OpenFile(cfg.Logger.Server, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	file2, err := os.OpenFile(cfg.Logger.Sqlite, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	file3, err := os.OpenFile(cfg.Logger.MW, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	var level slog.Level
	switch cfg.Logger.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	server := slog.New(slog.NewTextHandler(file1, &slog.HandlerOptions{
		Level: level,
	}))
	sqlite := slog.New(slog.NewTextHandler(file2, &slog.HandlerOptions{
		Level: level,
	}))
	mw := slog.New(slog.NewTextHandler(file3, &slog.HandlerOptions{
		Level: level,
	}))

	return &Logger{Server: server, Sqlite: sqlite, MW: mw}, nil
}
