package api

import(
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/models"
	"movie_searcher/middlewares"
)

func FetchSimilarMovies() echo.HandlerFunc {
	return func(c echo.Context) error {
		dbs := c.Get("dbs").(*middlewares.DatabaseClient)
		movies := []models.Movie{}
		// dbs.DB.Find(&movies)
		dbs.DB.Debug().Select("title").Find(&movies)
		return c.JSON(fasthttp.StatusOK, movies)
	}
}
