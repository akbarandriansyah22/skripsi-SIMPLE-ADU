package dto

type AIResponse struct {
	CleanedText string   `json:"cleaned_text"`
	Tokens      []string `json:"tokens"`
	Score       int      `json:"score"`
	Sentimen    string   `json:"sentiment"`
	Urgensi     string   `json:"urgency"`
}
