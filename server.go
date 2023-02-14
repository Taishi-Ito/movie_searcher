package main

import(
	"movie_searcher/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"net/http"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
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
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	routes.Init(e)
	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":8081"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!!!!")
}
