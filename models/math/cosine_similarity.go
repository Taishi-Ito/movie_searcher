package math

import (
	"gonum.org/v1/gonum/mat"
)

func CalcCosineSimilarity(input_vec, compared_vec []float64) float64 {
	input_vec_dense := mat.NewVecDense(768, input_vec)
	compared_vec_dense := mat.NewVecDense(768, compared_vec)
	dot := mat.Dot(input_vec_dense, compared_vec_dense)
	input_vec_norm := mat.Norm(input_vec_dense, 2)
	compared_vec_norm := mat.Norm(compared_vec_dense, 2)
	cosine_similarity := dot / (input_vec_norm * compared_vec_norm)
	return cosine_similarity
}
