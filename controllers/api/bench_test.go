package api

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"movie_searcher/databases"
	"movie_searcher/middlewares"
	"movie_searcher/models/movie"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkFetchSimilarMovies(b *testing.B) {
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

	b.ResetTimer()
	FetchSimilarMovies()(c)
	b.StopTimer()

	ordered_top_movies := []movie.Movie{}
	err := json.Unmarshal(rec.Body.Bytes(), &ordered_top_movies)

	if err != nil {
		b.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(ordered_top_movies) != 10 {
		b.Fatalf("unexpected number of movies: got %d, want 10", len(ordered_top_movies))
	}
}
