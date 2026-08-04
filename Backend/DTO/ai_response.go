package dto

type AIResponse struct {
	OriginalText         string           `json:"original_text"`
	CleanedText          string           `json:"cleaned_text"`
	Tokens               []string         `json:"tokens"`
	Score                int              `json:"score"`
	Sentimen             string           `json:"sentiment"`
	Urgensi              string           `json:"urgency"`
	PositiveScore        int              `json:"positive_score"`
	NegativeScore        int              `json:"negative_score"`
	SentimentScore       int              `json:"sentiment_score"`
	SentimentLabel       string           `json:"sentiment_label"`
	MatchedWords         []map[string]any `json:"matched_words"`
	SentimentExplanation string           `json:"sentiment_explanation"`
	UrgencyScore         int              `json:"urgency_score"`
	UrgencyLabel         string           `json:"urgency_label"`
	UrgencyReason        string           `json:"urgency_reason"`
}
