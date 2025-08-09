package app

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/kkvaleriy/istokAuth/internal/auth/delivery/http/v1"
	httperrors "github.com/kkvaleriy/istokAuth/internal/auth/delivery/http/v1/errors"
	"github.com/kkvaleriy/istokAuth/internal/auth/repository/postgres"
	"github.com/kkvaleriy/istokAuth/internal/auth/usecase"
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

type tokenParams interface {
	SecretKey() string
	RefreshTTL() time.Duration
	AccessTTL() time.Duration
}

type server struct {
	echo *echo.Echo
	cfg  serverCfg
}

type app struct {
	db          *pgxpool.Pool
	server      *server
	tokenParams tokenParams
	log         logger
}

func New(db *pgxpool.Pool, e *echo.Echo, tParams tokenParams, sCfg serverCfg, log logger) *app {
	return &app{
		db:          db,
		server:      &server{echo: e, cfg: sCfg},
		tokenParams: tParams,
		log:         log,
	}
}

func (app *app) Run() error {

	sysQuit := make(chan os.Signal, 1)

	signal.Notify(sysQuit, syscall.SIGTERM, syscall.SIGKILL)
	gracefulShutdownWG := &sync.WaitGroup{}
	gracefulShutdownWG.Add(1)

	app.server.echo.HTTPErrorHandler = httperrors.ErrorsHandler(app.log)

	domain := app.server.echo.Group("/api/v1")
	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	pg := postgres.New(app.db, app.log)
	uc := usecase.NewUserService(app.tokenParams, pg, app.log)
	ht := v1.NewHandler(uc, app.log)

	ht.Routes(domain, app.tokenParams.SecretKey())

	err := app.server.echo.Start(app.server.cfg.ServerPort())
	if err != nil {
		app.log.Fatal("server startup error", "error", err.Error())
	}

	go func() {
		<-sysQuit
		app.db.Close()
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
