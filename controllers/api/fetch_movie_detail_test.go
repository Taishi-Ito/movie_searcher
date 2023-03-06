package api

import (
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

func TestFetchMovieDetail(t *testing.T) {
	godotenv.Load()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/show/:id", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("119")

	session, _ := database.Connect()
	d := middlewares.DatabaseClient{DB: session}
	defer d.DB.Close()
	d.DB.LogMode(true)
	c.Set("dbs", &d)


	response := FetchMovieDetail()(c)
	if response != nil {
		t.Errorf("failed to fetch movies detail: %v", response)
		return
	}

	movie := movie.Movie{}
	err := json.Unmarshal(rec.Body.Bytes(), &movie)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
		return
	}

	if !(movie.Title == "グリーンマイル") {
		t.Errorf("Expected to get グリーンマイル, but got '%s' instead", movie.Title)
		return
	}
	if !(movie.Year == 1999) {
		t.Errorf("Expected to get Year, but got '%d' instead", movie.Year)
		return
	}
	if !(len(movie.Summary) > 0) {
		t.Errorf("Expected to get Summary, but got '%s' instead", movie.Summary)
		return
	}
	if !(movie.ImdbRank == 27) {
		t.Errorf("Expected to get ImdbRank, but got '%d' instead", movie.ImdbRank)
		return
	}
	if !(movie.PrimeRating == 4.6) {
		t.Errorf("Expected to get PrimeRating, but got '%f' instead", movie.PrimeRating)
		return
	}
	if !(movie.ImdbRating == 8.6) {
		t.Errorf("Expected to get ImdbRating, but got '%f' instead", movie.ImdbRating)
		return
	}
	if !(movie.AverageRating == 6.6) {
		t.Errorf("Expected to get AverageRating, but got '%f' instead", movie.AverageRating)
		return
	}
	if !(movie.PrimeReviewNum == 1279) {
		t.Errorf("Expected to get PrimeReviewNum, but got '%d' instead", movie.PrimeReviewNum)
		return
	}
	if !(movie.FilmLength == 188) {
		t.Errorf("Expected to get FilmLength, but got '%d' instead", movie.FilmLength)
		return
	}
}
