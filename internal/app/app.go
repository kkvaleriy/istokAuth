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
	"github.com/labstack/gommon/log"
)

type app struct {
	db   *pgxpool.Pool
	echo *echo.Echo
}

func New(db *pgxpool.Pool, e *echo.Echo) *app {
	return &app{
		db:   db,
		echo: e,
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

	pg := postgres.New(app.db)
	uc := usecase.New(pg)
	ht := v1.NewHandler(uc)

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
