package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	handlers "github.com/ykode/srp_demo/server/internal/handlers"
)

func main() {
	log.Info("Starting SRP...")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	identityHandlers := handlers.NewIdentityHandler()

	idGroup := e.Group("identities")
	identityHandlers.Mount(idGroup)

	e.Logger.Fatal(e.Start(":4000"))
}
