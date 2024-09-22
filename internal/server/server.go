package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"gofemart/internal/logger"
	"net/http"
)

type Server struct {
	S   *http.Server
	ctx context.Context
}

// NewServer Создаём сервер приложения
func NewServer(ctx context.Context, router chi.Router, address string) *Server {
	logger.Log.Infof("Running server on %s", address)
	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	return &Server{
		S:   &server,
		ctx: ctx,
	}
}

// Close закрытие сервера
func (s *Server) Close() {
	// Заставляем завершиться сервер и ждём его завершения
	err := s.S.Shutdown(s.ctx)
	if err != nil {
		logger.Log.Errorf("Failed to shutdown server: %v", err)
	}
	logger.Log.Info("Server stop")
}
