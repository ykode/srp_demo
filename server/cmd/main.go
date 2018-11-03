package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"os"

	handlers "github.com/ykode/srp_demo/server/internal/handlers"
	storage "github.com/ykode/srp_demo/server/internal/storage/postgresql"
)

func main() {
	log.Info("Starting SRP...")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	s, err := storage.NewPostgreSQLStorage("localhost", "postgres", "", "srp_demo", 5432, storage.SSLDisable)

	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	idGroup := e.Group("identities")
	identityHandlers := handlers.NewIdentityHandler(s, s)
	identityHandlers.Mount(idGroup)

	sessionGroup := e.Group("sessions")
	sessionHandlers := handlers.NewSessionHandler(s, s, s)
	sessionHandlers.Mount(sessionGroup)

	e.Logger.Fatal(e.Start(":4000"))
}
