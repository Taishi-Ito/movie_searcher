package routes

import (
	"github.com/labstack/echo/v4"
	"movie_searcher/controllers/api"
)

func Init(e *echo.Echo) {
	g := e.Group("/api")
	{
		g.POST("/similar", api.FetchSimilarMovies())
		g.GET("/show/:id", api.FetchMovieDetail())
	}
}
