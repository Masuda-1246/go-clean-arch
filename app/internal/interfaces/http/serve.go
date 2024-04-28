package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Masuda-1246/go-clean-arch/internal/interfaces/openapi/openapi"
)

type Server struct{}

// Login implements the login operation.
func (s *Server) Login(ctx echo.Context) error {
	// 実装するロジック
	return ctx.String(http.StatusOK, "Login successful")
}

// CheckHealthy implements the health check operation.
func (s *Server) CheckHealthy(ctx echo.Context) error {
	// 実装するロジック
	return ctx.String(http.StatusOK, "Everything is healthy")
}

func Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	server := &Server{}

	openapi.RegisterHandlers(e, server)

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := e.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
}
