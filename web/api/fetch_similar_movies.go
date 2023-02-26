package api

import(
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/models"
	"movie_searcher/middlewares"
	"movie_searcher/web/vender"
	"movie_searcher/web/utils/calculation"
	"encoding/json"
	"sort"
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

		// sentence-vector-generatorにリクエストを送信する
		input_vec := vender.FetchSentenceVector(request.Text)

		// DBからMovieの全データを取得する
		dbs := c.Get("dbs").(*middlewares.DatabaseClient)
		movies := []models.Movie{}
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
			cosine_similarity := calculation.CalcCosineSimilarity(input_vec, compared_vec)
			rankings = append(rankings, IdSimilarity{Id: movie.ID, Similarity: cosine_similarity})
		}

		// 類似度トップ10のデータを取得
		sort.Slice(rankings, func(i, j int) bool {return rankings[i].Similarity > rankings[j].Similarity})
		topMoviesID := make([]uint, 0)
		for _, ranking := range rankings[:10] {
            topMoviesID = append(topMoviesID, ranking.Id)
        }
		top_movies := []models.Movie{}
		dbs.DB.Debug().Select([]string{"id", "title", "year",}).Where(topMoviesID).Find(&top_movies)

		// 取得したデータを類似度順に並び替え
		ordered_top_movies := []models.Movie{}
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
