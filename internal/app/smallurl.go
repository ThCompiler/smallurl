package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"smallurl/internal/app/config"
	httpv1 "smallurl/internal/app/delivery/http/v1"
	"smallurl/internal/pkg/shortcut"
	grpcv1 "smallurl/internal/shortcut/delivery/grpc/v1"
	sh "smallurl/internal/shortcut/delivery/http/v1/handlers"
	"smallurl/internal/shortcut/repository/memory"
	"smallurl/internal/shortcut/repository/postgres"
	su "smallurl/internal/shortcut/usecase"
	generatedv1 "smallurl/pkg/grpc/v1"
	"smallurl/pkg/logger"
	"smallurl/pkg/server"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func initPgPool(cfg *config.Config, l logger.Interface) *pgxpool.Pool {
	// Postgres
	cfx, err := pgxpool.ParseConfig(cfg.Postgres.URL)
	if err != nil {
		l.Fatal("[App] Init - postgres.New: %s", err)
	}

	cfx.MaxConns = int32(cfg.Postgres.MaxConnections)
	cfx.MinConns = int32(cfg.Postgres.MinConnections)
	cfx.MaxConnIdleTime = time.Duration(cfg.Postgres.TTLIDleConnections) * time.Millisecond

	pg, err := pgxpool.NewWithConfig(context.Background(), cfx)
	if err != nil {
		l.Fatal("[App] Init - postgres.New: %s", err)
	}

	if err = pg.Ping(context.Background()); err != nil {
		pg.Close()
		l.Fatal("[App] Init - can't check connection to sql with error %s", err)
	}

	l.Info("[App] Init - success check connection to postgresql")

	return pg
}

// Run -
// nolint: revive // Точка входа, которую сложно ещё сильнее разбить на подфункции, чтобы сократить размер
func Run(cfg *config.Config) {
	// Logger
	l, logFile := prepareLogger(cfg.LoggerInfo)

	defer func() {
		_ = l.Sync() // nolint: errcheck //нет смысла логировать ошибку, при выключении сервер

		if logFile != nil {
			_ = logFile.Close() // nolint: errcheck // нет смысла логировать ошибку закрытия лога,
			// при выключении сервера
		}
	}()

	var rep su.Repository
	if cfg.UseInMemory {
		rep = memory.NewRepository()
	} else {
		// Databases
		pg := initPgPool(cfg, l)
		defer pg.Close()

		// Repository
		rep = postgres.NewRepository(pg)
	}

	// Usecase
	usc := su.NewUsecase(rep, func() su.Shortcut { return shortcut.NewHashShortcut() })

	// Handlers
	shortcutHandlers := sh.NewShortcutHandlers(usc)

	// Routes
	routes := prepareRoutes(shortcutHandlers)

	router, err := httpv1.NewRouter("/api", routes, cfg.Mode, l)
	if err != nil {
		l.Fatal("[App] Init - init handler error: %s", err)
	}

	httpServer := server.New(router, server.Port(cfg.HTTP.Port))

	// GRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GRPC.Hostname, cfg.GRPC.Port))
	if err != nil {
		l.Error("failed to listen: %v", err)

		return
	}

	grpcServer := grpc.NewServer()
	generatedv1.RegisterShortcutServer(grpcServer, grpcv1.NewShortcutService(usc, l))

	grpcNotify := make(chan error)
	go func() {
		grpcNotify <- grpcServer.Serve(lis)
		close(grpcNotify)
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	l.Info("[App] Start - server started")

	select {
	case s := <-interrupt:
		l.Info("[App] Run - signal: " + s.String())
	case err := <-grpcNotify:
		l.Error(fmt.Errorf("[App] Run - grpcServer.Notify: %w", err))
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("[App] Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	grpcServer.Stop()

	err = httpServer.Shutdown()
	if err != nil {
		l.Fatal(fmt.Errorf("[App] Stop - httpServer.Shutdown: %w", err))
	}

	l.Info("[App] Stop - server stopped")
}
