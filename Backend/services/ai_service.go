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
		Score:    prediction.Score,
		Sentimen: prediction.Sentiment,
		Urgensi:  prediction.Urgency,
	}

	return &result, nil
}
