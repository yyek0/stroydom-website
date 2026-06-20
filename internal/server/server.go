package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yyek0/stroydom-website/internal/handler"
)

type Server struct {
	httpHandlers *handler.Handlers
}

func NewServer(httpHandler *handler.Handlers) *Server {
	return &Server{
		httpHandlers: httpHandler,
	}
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()

	router.Path("/health").Methods("GET").HandlerFunc(s.httpHandlers.HandleCheckHealth)
	router.Path("/leads").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateLead)
	router.Path("/leads").Methods("GET").Queries("id", "{id}").HandlerFunc(s.httpHandlers.HandleGetLead)
	router.Path("/leads").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllLeads)
	router.Path("/leads").Methods("DELETE").Queries("id", "{id}").HandlerFunc(s.httpHandlers.HandleDeleteLead)

	return http.ListenAndServe(":8080", router)
}
