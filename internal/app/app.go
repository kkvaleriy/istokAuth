package app

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/kkvaleriy/istokAuthorization/internal/authorization/delivery/http/v1"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/repository/postgres"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/usecase"
	"github.com/labstack/echo/v4"
)

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type serverCfg interface {
	ServerPort() string
}

type server struct {
	echo *echo.Echo
	cfg  serverCfg
}

type app struct {
	db     *pgxpool.Pool
	server *server
	log    logger
}

func New(db *pgxpool.Pool, e *echo.Echo, sCfg serverCfg, log logger) *app {
	return &app{
		db:     db,
		server: &server{echo: e, cfg: sCfg},
		log:    log,
	}
}

func (app *app) Run() error {

	sysQuit := make(chan os.Signal, 1)

	signal.Notify(sysQuit, syscall.SIGTERM, syscall.SIGKILL)
	gracefulShutdownWG := &sync.WaitGroup{}
	gracefulShutdownWG.Add(1)

	domain := app.server.echo.Group("/api/v1")
	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	pg := postgres.New(app.db, app.log)
	uc := usecase.New(pg, app.log)
	ht := v1.NewHandler(uc, app.log)

	ht.Routes(domain)

	err := app.server.echo.Start(app.server.cfg.ServerPort())
	if err != nil {
		app.log.Fatal("server startup error", "error", err.Error())
	}

	go func() {
		<-sysQuit
		app.log.Info("the database connection has been successfully closed")
		err = app.server.echo.Close()
		if err != nil {
			app.log.Fatal("the server shutdown error", err.Error())
		}
		app.log.Info("the server has been successfully stopped")
	}()

	gracefulShutdownWG.Wait()
	return nil
}
