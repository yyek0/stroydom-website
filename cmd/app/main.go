package main

import (
	"context"
	"os"

	"github.com/joho/godotenv" // <-- вот эта строка нужна

	"github.com/yyek0/stroydom-website/internal/database"
	"github.com/yyek0/stroydom-website/internal/handler"
	"github.com/yyek0/stroydom-website/internal/server"
)

func main() {
	ctx := context.Background()

	// postgresql://username:password@host:port/database_name

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// 2. Достаем строку подключения
	connString := os.Getenv("DB_CONN")
	if connString == "" {
		panic("conn string is empty")
	}

	PostgresDB, err := database.NewDatabase(ctx, "postgres://postgres:1@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}

	if err := PostgresDB.Init(ctx); err != nil {
		panic(err)
	}

	handlers := handler.NewHandler(PostgresDB)

	serv := server.NewServer(handlers)

	if err := serv.StartServer(); err != nil {
		panic(err)
	}

}
