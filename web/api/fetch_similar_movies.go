package api

import(
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
	"movie_searcher/models"
	"movie_searcher/middlewares"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"encoding/json"
	"net/http"
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"sort"
)

type Response struct {
	EncodedText string  `json:"encoded_text"`
}

func FetchSimilarMovies() echo.HandlerFunc {
	return func(c echo.Context) error {
		text := c.QueryParam("text")

		// sentence-vector-generatorにリクエストを送信する
		jsonStr := []byte(fmt.Sprintf(`{"text": "%s"}`, text))
		req, err := http.NewRequest("POST",
		"http://localhost:8000/generate",
		bytes.NewBuffer([]byte(jsonStr)),
		)
		if err != nil {
			logrus.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		response := Response{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			logrus.Fatal(err)
		}
		input_vec := []float64{}
		json.Unmarshal([]byte(response.EncodedText), &input_vec)
		defer resp.Body.Close()


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
			if len(rankings) == 11 {
				sort.Slice(rankings, func(i, j int) bool {return rankings[i].Similarity > rankings[j].Similarity})
				rankings = rankings[:10]
			}
		}

		// 類似度トップ10のデータを取得
		top_movies := []models.Movie{}
		dbs.DB.Debug().Select([]string{"id", "title", "year",}).Where([]uint{rankings[0].Id, rankings[1].Id, rankings[2].Id, rankings[3].Id, rankings[4].Id, rankings[5].Id, rankings[6].Id, rankings[7].Id, rankings[8].Id, rankings[9].Id}).Find(&top_movies)


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
