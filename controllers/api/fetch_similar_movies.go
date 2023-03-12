package api

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/middlewares"
	"movie_searcher/models/math"
	"movie_searcher/models/movie"
	"movie_searcher/models/nlp"
	"sort"
	"sync"
)

func getAllMovies(dbs *middlewares.DatabaseClient, wg *sync.WaitGroup, ch chan []movie.Movie) {
	defer wg.Done()
	movies := []movie.Movie{}
	dbs.DB.Debug().Select([]string{"id", "average_vector"}).Find(&movies)
	ch <- movies
}

func FetchSimilarMovies() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request struct {
			Text string `json:"text"`
		}
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(2)

		// 入力文を文ベクトルに変換する
		ch_vec := make(chan []float64)
		go nlp.FetchSentenceVector(request.Text, &wg, ch_vec)
		input_vec := <-ch_vec

		// DBからMovieの全データを取得する
		ch_db := make(chan []movie.Movie)
		dbs := c.Get("dbs").(*middlewares.DatabaseClient)
		go getAllMovies(dbs, &wg, ch_db)
		movies := <-ch_db

		wg.Wait()

		// ベクトルの類似度を計算する
		type IdSimilarity struct {
			Id         uint
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
		sort.Slice(rankings, func(i, j int) bool { return rankings[i].Similarity > rankings[j].Similarity })
		top_movies_ids := make([]uint, 0)
		for _, ranking := range rankings[:10] {
			top_movies_ids = append(top_movies_ids, ranking.Id)
		}
		top_movies := []movie.Movie{}
		dbs.DB.Debug().Select([]string{"id", "title", "year"}).Where(top_movies_ids).Find(&top_movies)

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
