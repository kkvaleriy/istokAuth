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

	domain := app.echo.Group("/api/v1")
	domain.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	app.log.Debug("added route for /api/v1")

	pg := postgres.New(app.db, app.log)
	uc := usecase.New(pg, app.log)
	ht := v1.NewHandler(uc, app.log)

	ht.Routes(domain)

	err := app.echo.Start(":8080")
	if err != nil {
		log.Errorf("error of start server: %v", err)
	}

	go func() {
		<-sysQuit

		err = app.echo.Close()
		if err != nil {
			log.Fatalf("can't stop server: %s", err.Error())
		}
	}()

	gracefulShutdownWG.Wait()
	return nil
}
