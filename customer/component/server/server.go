package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DuongVu089x/interview/customer/component/appctx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
	port   string
	appCtx appctx.AppContext
}

func NewServer(e *echo.Echo, port string, appCtx appctx.AppContext) *Server {
	return &Server{
		echo:   e,
		port:   port,
		appCtx: appCtx,
	}
}

func (s *Server) Start() error {

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", s.port)
		log.Printf("Starting server on %s", addr)
		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown services
	if err := s.shutdown(ctx); err != nil {
		s.appCtx.GetKafkaConsumer().Close()
		return fmt.Errorf("error during shutdown: %v", err)
	}

	return nil
}

func (s *Server) shutdown(ctx context.Context) error {

	// Shutdown HTTP server
	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("error during server shutdown: %v", err)
	}

	return nil
}
