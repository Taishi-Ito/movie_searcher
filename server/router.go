package routes

import (
	"movie_searcher/controllers/api"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	g := e.Group("/api")
	{
		g.POST("/similar", api.FetchSimilarMovies())
		g.GET("/show/:id", api.FetchMovieDetail())
	}
}
