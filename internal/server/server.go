package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/yyek0/stroydom-website/internal/handler"
	"go.uber.org/zap"
)

type Server struct {
	httpHandlers *handler.Handlers
}

func NewServer(httpHandler *handler.Handlers) *Server {
	return &Server{
		httpHandlers: httpHandler,
	}
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func logRequestDetails(logger *zap.Logger, start time.Time, r *http.Request) {
	logger.Info("Входящий HTTP-запрос",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("ip", r.RemoteAddr),
		zap.Duration("duration", time.Since(start)),
	)
}

func handlePanic(logger *zap.Logger, w http.ResponseWriter, panicErr interface{}) {
	logger.Error("КРИТИЧЕСКАЯ ОШИБКА (PANIC)",
		zap.Any("error", panicErr),
		zap.String("stack", string(debug.Stack())),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"status":"error","message":"Внутренняя ошибка сервера"}`))
}

func (s *Server) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Вызываем наш хелпер, если случилась беда
				handlePanic(s.httpHandlers.Logger, w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Server) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logRequestDetails(s.httpHandlers.Logger, start, r)
	})
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()

	router.Use(s.RecoveryMiddleware)
	router.Use(s.CORSMiddleware)
	router.Use(s.LoggingMiddleware)

	router.Path("/api/health").Methods("GET").HandlerFunc(s.httpHandlers.HandleCheckHealth)
	router.Path("/api/leads").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateLead)
	router.Path("/api/leads").Methods("GET").Queries("id", "{id}").HandlerFunc(s.httpHandlers.HandleGetLead)
	router.Path("/api/leads").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllLeads)
	router.Path("/api/leads").Methods("DELETE").Queries("id", "{id}").HandlerFunc(s.httpHandlers.HandleDeleteLead)

	srv := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.httpHandlers.Logger.Fatal("Критическая ошибка при запуске сервера", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	sig := <-quit

	s.httpHandlers.Logger.Info("Получен сигнал остановки сервера. Инициирован Graceful Shutdown",
		zap.String("signal", sig.String()),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.httpHandlers.Logger.Error("Сервер завершил работу с ошибкой таймаута", zap.Error(err))
		return err
	}

	s.httpHandlers.Logger.Info("Сервер успешно и безопасно остановлен")
	return nil
}
