package main

import(
	"movie_searcher/routes"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"net/http"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	routes.Init(e)
	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":8081"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!!!!")
}
