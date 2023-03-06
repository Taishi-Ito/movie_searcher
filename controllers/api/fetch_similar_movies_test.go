package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"movie_searcher/middlewares"
	"github.com/labstack/echo/v4"
	"encoding/json"
	"movie_searcher/models/movie"
	"github.com/joho/godotenv"
	"movie_searcher/databases"
)

func TestFetchSimilarMovies(t *testing.T) {
	godotenv.Load()
	e := echo.New()
	jsonStr := `{"text":"感動してみんなで泣ける映画"}`
	req := httptest.NewRequest(http.MethodPost, "/api/similar", bytes.NewReader([]byte(jsonStr)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	session, _ := database.Connect()
	d := middlewares.DatabaseClient{DB: session}
	defer d.DB.Close()
	d.DB.LogMode(true)
	c.Set("dbs", &d)


	response := FetchSimilarMovies()(c)
	if response != nil {
		t.Errorf("failed to fetch similar movies: %v", response)
		return
	}

	ordered_top_movies := []movie.Movie{}
	err := json.Unmarshal(rec.Body.Bytes(), &ordered_top_movies)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(ordered_top_movies) != 10 {
		t.Errorf("Expected slice length to be 10, but got '%d'", len(ordered_top_movies))
		return
	}
	for _, v := range ordered_top_movies {
		if !(v.CosineSimilarity > 0) {
			t.Errorf("Expected to get CosineSimilarity, but got '%f' instead", v.CosineSimilarity)
			break
		}
		if !(len(v.Title) > 0) {
			t.Errorf("Expected to get title, but got '%s' instead", v.Title)
			break
		}
		if !(v.Year >= 0) {
			t.Errorf("Expected to get Year, but got '%d' instead", v.Year)
			break
		}
	}
}
