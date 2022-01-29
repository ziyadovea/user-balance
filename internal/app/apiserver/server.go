package apiserver

import (
	"context"
	"net/http"
	"time"
)

// Server - структура для представления сервера приложения
type Server struct {
	httpServer *http.Server
}

// NewServer - конструктор для сервера
func NewServer(port string, handler http.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}}
}

// Run запускат сервер
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown реализует gracefully shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
