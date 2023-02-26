package routes

import (
	"movie_searcher/web/api"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	g := e.Group("/api")
	{
		g.POST("/similar", api.FetchSimilarMovies())
	}
}
