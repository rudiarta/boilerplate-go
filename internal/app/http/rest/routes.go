package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rudiarta/boilerplate-go/internal/app/config"
	"github.com/rudiarta/boilerplate-go/internal/app/registry"
	"github.com/rudiarta/boilerplate-go/pkg/goroutinecheck"

	"github.com/google/uuid"

	v1handler "github.com/rudiarta/boilerplate-go/internal/app/http/rest/v1"
)

var appCtx app

type app struct {
	repo    registry.RepositoryRegistry
	usecase registry.UsecaseRegistry
	cfg     config.ConfigCtx
}

func initApp() app {
	cfg := config.NewConfigCtx()
	repo := registry.NewRepositoryRegistry(cfg)
	return app{
		repo:    repo,
		cfg:     cfg,
		usecase: registry.NewUsecaseRegistry(repo),
	}
}

func Run() {
	appCtx = initApp()

	e := echo.New()

	// register handler
	v1(e)

	// Start server
	go func() {
		keyGoRoutine := uuid.New().String() //Generate unique identifier
		goroutinecheck.IncreaseTotalGoRoutineCount(keyGoRoutine)
		defer goroutinecheck.DecreaseTotalGoRoutineCount(keyGoRoutine)

		envValue, _ := appCtx.cfg.GetEnvironmentValue()
		portHost := fmt.Sprintf(":%s", envValue.Port)
		if err := e.Start(portHost); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 25 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// do something below this line
	ticker := time.NewTicker(1 * time.Second)
	log.Printf(`total go routine: %d `, goroutinecheck.GetTotalGoRoutineCount())

loop:
	for {
		select {
		case <-ctx.Done():

			log.Print("context time out, force killing")
			log.Printf("with total goroutine left: %d", goroutinecheck.GetTotalGoRoutineCount())

			break loop

		case <-ticker.C:
			log.Printf(`total go routine: %d `, goroutinecheck.GetTotalGoRoutineCount())
			if goroutinecheck.GetTotalGoRoutineCount() <= 0 {
				break loop
			}
			continue

		}
	}
}

func v1(e *echo.Echo) {

	v1 := e.Group("/v1")

	// AccountHandler
	accountHandler := v1handler.NewAccountHandler(appCtx.cfg, appCtx.usecase.AccountUsecase())
	v1account := v1.Group("/account")
	v1account.GET("/create", accountHandler.CreateAccount)
}
