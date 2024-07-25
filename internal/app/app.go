package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"

	"boilerplate/internal/config"
	"boilerplate/router"
)

type App struct {
	DB     *gorm.DB
	Router *router.Router
	Server *http.Server
}

func (a *App) Start(ctx context.Context, conf *config.Config) {
	a.Router = &router.Router{DB: a.DB}
	err := a.Router.Routes(conf)
	if err != nil {
		panic(err)
	}

	a.Server = &http.Server{Addr: ":3000"}

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	log.Println("Server started at :3000")

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctxShutDown); err != nil {
		fmt.Printf("server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server exited properly")
	return
}

func (a *App) Shutdown() {
	fmt.Println("App shutdown")
}
