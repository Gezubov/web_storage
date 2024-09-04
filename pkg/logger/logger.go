package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	// Настройка логгера: выводим в консоль с уровнем INFO
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Создаем глобальный логгер
	Logger = slog.New(handler)
}
