package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	Score     int    `json:"score"`
	Sentiment string `json:"sentiment"`
	Urgency   string `json:"urgency"`
}

func NewAIService() *AIService {
	baseURL := os.Getenv("AI_SERVICE")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	return &AIService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *AIService) Analyze(req dto.AIRequest) (*dto.AIResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(s.baseURL+"/predict", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("AI service mengembalikan status %d", resp.StatusCode)
	}

	var prediction aiPredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, err
	}

	if prediction.Sentiment == "" || prediction.Urgency == "" {
		return nil, errors.New("response AI tidak lengkap")
	}

	result := dto.AIResponse{
		Score:      prediction.Score,
		Sentimen:   prediction.Sentiment,
		Urgensi:    prediction.Urgency,
		Confidence: float64(prediction.Score),
	}

	return &result, nil
}
