package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dto "backend/DTO"
)

func newTestAIService(t *testing.T, handler http.HandlerFunc) *AIService {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &AIService{
		baseURL: server.URL,
		client:  server.Client(),
	}
}

func TestAIServiceAnalyzeValidResponse(t *testing.T) {
	service := newTestAIService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("content-type = %s", r.Header.Get("Content-Type"))
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if _, ok := payload["deskripsi"]; !ok {
			t.Fatal("request tidak mengirim field deskripsi")
		}
		if _, ok := payload["judul"]; ok {
			t.Fatal("request masih mengirim field judul")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"score":0,"sentiment":"Netral","urgency":"Rendah"}`))
	})

	result, err := service.Analyze(dto.AIRequest{Deskripsi: "Wifi lambat"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Score != 0 || result.Sentimen != "Netral" || result.Urgensi != "Rendah" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestAIServiceAnalyzeRejectsErrorStatus(t *testing.T) {
	for _, statusCode := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			service := newTestAIService(t, func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "bad ai response", statusCode)
			})

			if _, err := service.Analyze(dto.AIRequest{Deskripsi: "Wifi lambat"}); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestAIServiceAnalyzeTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()

	service := &AIService{
		baseURL: server.URL,
		client:  &http.Client{Timeout: 10 * time.Millisecond},
	}

	if _, err := service.Analyze(dto.AIRequest{Deskripsi: "Wifi lambat"}); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestAIServiceAnalyzeRejectsMalformedJSON(t *testing.T) {
	service := newTestAIService(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	})

	if _, err := service.Analyze(dto.AIRequest{Deskripsi: "Wifi lambat"}); err == nil {
		t.Fatal("expected malformed json error")
	}
}

func TestAIServiceAnalyzeRequiresSentimentAndUrgency(t *testing.T) {
	tests := map[string]string{
		"missing sentiment": `{"score":1,"urgency":"Rendah"}`,
		"missing urgency":   `{"score":1,"sentiment":"Netral"}`,
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			service := newTestAIService(t, func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(body))
			})

			if _, err := service.Analyze(dto.AIRequest{Deskripsi: "Wifi lambat"}); err == nil {
				t.Fatal("expected incomplete response error")
			}
		})
	}
}
