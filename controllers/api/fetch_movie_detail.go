package api

import(
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/models/movie"
	"movie_searcher/middlewares"
)

func FetchMovieDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		dbs := c.Get("dbs").(*middlewares.DatabaseClient)
		movie := movie.Movie{}
		dbs.DB.Debug().Select([]string{"title", "year", "summary", "imdb_rank", "prime_rating", "imdb_rating", "average_rating", "prime_id", "prime_review_num", "film_length"}).Where("id = ?", id).First(&movie)
		return c.JSON(fasthttp.StatusOK, movie)
	}
}
