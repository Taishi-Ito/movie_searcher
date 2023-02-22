package models

type Movie struct {
	ID uint `json:"id" gorm:"primary_key"`
	Title string `json: "title"`
	RawTitle string `json: "raw_title"`
	Year uint `json: "year"`
	Summary string `json: "summary"`
	ImdbRank uint `json: "imdb_rank"`
	PrimeRating float32 `json: "prime_rating"`
	ImdbRating float32 `json: "imdb_rating"`
	AverageRating float32 `json: "average_rating"`
	PrimeId string `json: "prime_id"`
	PrimeReviewNum string `json: "prime_review_num"`
	FilmLength string `json: "film_length"`
	AverageVector string `json: "average_vector"`
	CosineSimilarity float64
}
