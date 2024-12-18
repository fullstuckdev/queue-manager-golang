package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"queue-manager/internal/api/handlers"
	"queue-manager/internal/api/routes"
	"queue-manager/internal/config"
)

type Server struct {
    echo   *echo.Echo
    logger *logrus.Logger
    cfg    *config.Config
}

func NewServer(cfg *config.Config, logger *logrus.Logger, tenantHandler *handlers.TenantHandler) *Server {
    e := echo.New()
    
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())
    e.Use(middleware.CORS())
    
    routes.SetupRoutes(e, tenantHandler)
    
    return &Server{
        echo:   e,
        logger: logger,
        cfg:    cfg,
    }
}

func (s *Server) Start() error {
    go func() {
        addr := fmt.Sprintf(":%s", s.cfg.Server.Port)
        if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
            s.logger.Fatalf("Failed to start server: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    return s.echo.Shutdown(ctx)
} 