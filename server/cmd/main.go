package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	handlers "github.com/ykode/srp_demo/server/internal/handlers"
	storage "github.com/ykode/srp_demo/server/internal/storage/inmemory"
)

func main() {
	log.Info("Starting SRP...")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	idStorage := storage.NewInMemoryIdentityStorage()
	sessionStorage := storage.NewInMemorySessionStorage()

	idGroup := e.Group("identities")
	identityHandlers := handlers.NewIdentityHandler(idStorage, idStorage)
	identityHandlers.Mount(idGroup)

	sessionGroup := e.Group("sessions")
	sessionHandlers := handlers.NewSessionHandler(sessionStorage, sessionStorage, idStorage)
	sessionHandlers.Mount(sessionGroup)

	e.Logger.Fatal(e.Start(":4000"))
}
