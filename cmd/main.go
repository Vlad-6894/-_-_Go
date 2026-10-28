package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	core_logger "security/internal/core/logger"
	core_postgres_pool "security/internal/core/repository/postgres/pool"
	core_http_middleware "security/internal/core/transport/http/middleware"
	core_http_server "security/internal/core/transport/http/server"
	postgres_repository "security/internal/features/repository/postgres"
	"security/internal/features/service"
	security_http_transport "security/internal/features/transport/http"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewLoggerConfigMust())
	if err != nil {
		fmt.Println("init application logger error!")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Init connection pool!")
	pool, err := core_postgres_pool.NewPostgresConnectionPool(ctx, core_postgres_pool.NewPostgresConnectionConfigMust())
	if err != nil {
		logger.Fatal("Failed to init connection pool!", zap.Error(err))
	}
	defer pool.Close()

	logger.Info("Init feature subscriptions!")
	dangerRepository := postgres_repository.NewDangersPostgresRepository(pool)
	dangerService := service.NewDangerService(dangerRepository)
	dangerTransportHTTP := security_http_transport.NewHTTPSecurityTransport(dangerService)

	logger.Info("Init HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewHTTPServerConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter()
	apiVersionRouter.RegisterRoutes(dangerTransportHTTP.Routes()...)

	httpServer.RegisterRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error:", zap.Error(err))
	}
}
