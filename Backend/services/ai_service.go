package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	dto "backend/DTO"
)

type AIService struct {
	baseURL string
	client  *http.Client
}

type aiPredictResponse struct {
	OriginalText         string           `json:"original_text"`
	CleanedText          string           `json:"cleaned_text"`
	Tokens               []string         `json:"tokens"`
	Score                int              `json:"score"`
	Sentiment            string           `json:"sentiment"`
	Urgency              string           `json:"urgency"`
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

func NewAIService() *AIService {
	baseURL := os.Getenv("AI_SERVICE")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	return &AIService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *AIService) Analyze(req dto.AIRequest) (*dto.AIResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, s.baseURL+"/predict", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		errorBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		msg := strings.TrimSpace(string(errorBody))
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return nil, fmt.Errorf("AI service mengembalikan status %d: %s", resp.StatusCode, msg)
	}

	var prediction aiPredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, err
	}

	if prediction.Sentiment == "" || prediction.Urgency == "" {
		return nil, errors.New("response AI tidak lengkap")
	}

	result := dto.AIResponse{
		OriginalText:         prediction.OriginalText,
		CleanedText:          prediction.CleanedText,
		Tokens:               prediction.Tokens,
		Score:                prediction.Score,
		Sentimen:             prediction.Sentiment,
		Urgensi:              prediction.Urgency,
		PositiveScore:        prediction.PositiveScore,
		NegativeScore:        prediction.NegativeScore,
		SentimentScore:       prediction.SentimentScore,
		SentimentLabel:       prediction.SentimentLabel,
		MatchedWords:         prediction.MatchedWords,
		SentimentExplanation: prediction.SentimentExplanation,
		UrgencyScore:         prediction.UrgencyScore,
		UrgencyLabel:         prediction.UrgencyLabel,
		UrgencyReason:        prediction.UrgencyReason,
	}
	if result.SentimentScore == 0 && result.Score != 0 {
		result.SentimentScore = result.Score
	}
	if result.SentimentLabel == "" {
		result.SentimentLabel = result.Sentimen
	}
	if result.UrgencyLabel == "" {
		result.UrgencyLabel = result.Urgensi
	}

	return &result, nil
}
