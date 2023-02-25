package api

import(
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/models"
	"movie_searcher/middlewares"
	"movie_searcher/web/sentence_vector_generator"
	"gonum.org/v1/gonum/mat"
	"encoding/json"
	"sort"
)

func FetchSimilarMovies() echo.HandlerFunc {
	return func(c echo.Context) error {
		text := c.QueryParam("text")

		// sentence-vector-generatorにリクエストを送信する
		input_vec := sentence_vector_generator.FetchSentenceVector(text)

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
			input_vec_dense := mat.NewVecDense(768, input_vec)
			compared_vec_dense := mat.NewVecDense(768, compared_vec)
			dot := mat.Dot(input_vec_dense, compared_vec_dense)
			input_vec_norm := mat.Norm(input_vec_dense, 2)
			compared_vec_norm := mat.Norm(compared_vec_dense, 2)
			cosine_similarity := dot / (input_vec_norm * compared_vec_norm)

			// 類似度順に並び替える
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
