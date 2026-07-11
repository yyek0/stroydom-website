package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/yyek0/stroydom-website/internal/database"
	"github.com/yyek0/stroydom-website/internal/handler"
	"github.com/yyek0/stroydom-website/internal/logger"
	"github.com/yyek0/stroydom-website/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	// Инициализируем наш кастомный логгер
	appLogger := logger.InitLogger()
	defer appLogger.Sync() // Сбрасываем буферы при выключении

	connString := os.Getenv("DB_CONN")
	db, err := database.NewDatabase(context.Background(), connString)
	if err != nil {
		appLogger.Fatal("Не удалось подключиться к БД", zap.Error(err))
	}

	// Прокидываем логгер в хендлеры
	myHandlers := handler.NewHandler(db, appLogger)
	myServer := server.NewServer(myHandlers)

	appLogger.Info("Сервер успешно запущен")
	if err := myServer.StartServer(); err != nil {
		appLogger.Fatal("Сервер упал", zap.Error(err))
	}
}
