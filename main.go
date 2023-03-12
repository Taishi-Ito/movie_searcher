package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"movie_searcher/middlewares"
	"movie_searcher/server"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env")
	}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	e := echo.New()
	env := os.Getenv("ENV")
	if env == "prod" {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(os.Getenv("URL"))
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/certs")
		e.Pre(middleware.HTTPSWWWRedirect())
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middlewares.DatabaseService())

	// Routes
	routes.Init(e)
	switch env {
	case "prod":
		e.Logger.Fatal(e.StartAutoTLS(":443"))
	default:
		defaultAddr := ":8081"
		e.Logger.Fatal(e.Start(defaultAddr))
	}
}
