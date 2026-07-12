package dto

type AIResponse struct {
	Score    int    `json:"score"`
	Sentimen string `json:"sentimen"`
	Urgensi  string `json:"urgensi"`
}
