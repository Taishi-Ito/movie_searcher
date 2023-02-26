package api

import(
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/models/math"
	"movie_searcher/models/movie"
	"movie_searcher/models/nlp"
	"movie_searcher/middlewares"
	"encoding/json"
	"sort"
	"fmt"
)

func FetchSimilarMovies() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request struct {
			Text string `json:"text"`
		}
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		// 入力文を文ベクトルに変換する
		input_vec := nlp.FetchSentenceVector(request.Text)

		// DBからMovieの全データを取得する
		dbs := c.Get("dbs").(*middlewares.DatabaseClient)
		movies := []movie.Movie{}
		dbs.DB.Debug().Select([]string{"id","average_vector"}).Find(&movies)

		// ベクトルの類似度を計算する
		type IdSimilarity struct {
			Id uint
			Similarity float64
		}
		rankings := []IdSimilarity{}
		for _, movie := range movies {
			compared_vec := []float64{}
			json.Unmarshal([]byte(movie.AverageVector), &compared_vec)
			cosine_similarity := math.CalcCosineSimilarity(input_vec, compared_vec)
			rankings = append(rankings, IdSimilarity{Id: movie.ID, Similarity: cosine_similarity})
		}

		// 類似度トップ10のデータを取得
		sort.Slice(rankings, func(i, j int) bool {return rankings[i].Similarity > rankings[j].Similarity})
		topMoviesID := make([]uint, 0)
		for _, ranking := range rankings[:10] {
            topMoviesID = append(topMoviesID, ranking.Id)
        }
		top_movies := []movie.Movie{}
		dbs.DB.Debug().Select([]string{"id", "title", "year",}).Where(topMoviesID).Find(&top_movies)

		// 取得したデータを類似度順に並び替え
		ordered_top_movies := []movie.Movie{}
		for _, ranking := range rankings {
			for _, top_movie := range top_movies {
				if ranking.Id == top_movie.ID {
					top_movie.CosineSimilarity = ranking.Similarity
					ordered_top_movies = append(ordered_top_movies, top_movie)
					break
				}
			}
		}

		return c.JSON(fasthttp.StatusOK, ordered_top_movies)
	}
}

func FetchMovieDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		dbs := c.Get("dbs").(*middlewares.DatabaseClient)
		movie := movie.Movie{}
		dbs.DB.Debug().Select([]string{"title", "year", "summary", "imdb_rank", "prime_rating", "imdb_rating", "average_rating", "prime_id", "prime_review_num", "film_length"}).Where("id = ?", id).First(&movie)
		fmt.Println(movie)
		return c.JSON(fasthttp.StatusOK, movie)
	}
}
