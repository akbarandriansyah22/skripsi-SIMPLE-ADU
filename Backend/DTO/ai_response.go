package dto

type AIResponse struct {
	KategoriPrediksi string  `json:"kategori_prediksi"`
	Score            int     `json:"score"`
	Sentimen         string  `json:"sentimen"`
	Urgensi          string  `json:"urgensi"`
	Confidence       float64 `json:"confidence"`
}
