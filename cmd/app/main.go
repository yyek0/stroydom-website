package main

import (
	"github.com/yyek0/stroydom-website/internal/handler"
	"github.com/yyek0/stroydom-website/internal/server"
)

func main() {

	dummy := handler.DummyDB{}
	handlers := handler.NewHandler(&dummy)

	serv := server.NewServer(handlers)

	err := serv.StartServer()
	if err != nil {
		panic(err)
	}

}
