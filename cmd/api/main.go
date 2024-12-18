package main

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"

    "queue-manager/internal/api"
    "queue-manager/internal/api/handlers"
    "queue-manager/internal/config"
    "queue-manager/internal/queue"
    "queue-manager/internal/repository"
)

func main() {
    logger := logrus.New()
    logger.SetOutput(&lumberjack.Logger{
        Filename:   "logs/app.log",
        MaxSize:    100, 
        MaxBackups: 3,
        MaxAge:     28, 
        Compress:   true,
    })
    logger.SetFormatter(&logrus.JSONFormatter{})

    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Fatalf("Failed to load config: %v", err)
    }

    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
        cfg.DB.User,
        cfg.DB.Password,
        cfg.DB.Host,
        cfg.DB.Port,
        cfg.DB.Name,
    )
    
    dbPool, err := pgxpool.Connect(context.Background(), dbURL)
    if err != nil {
        logger.Fatalf("Unable to connect to database: %v", err)
    }
    defer dbPool.Close()

    rmq, err := queue.NewRabbitMQ(
        cfg.RabbitMQ.User,
        cfg.RabbitMQ.Password,
        cfg.RabbitMQ.Host,
        cfg.RabbitMQ.Port,
        logger,
    )
    if err != nil {
        logger.Fatalf("Failed to initialize RabbitMQ: %v", err)
    }

    tenantRepo := repository.NewTenantRepository(dbPool)

    tenantHandler := handlers.NewTenantHandler(tenantRepo, rmq, logger)

    server := api.NewServer(cfg, logger, tenantHandler)
    if err := server.Start(); err != nil {
        logger.Fatalf("Server error: %v", err)
        os.Exit(1)
    }
} 