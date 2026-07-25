package services

import (
	"bytes"
	"testing"

	dto "backend/DTO"
)

func TestHasilAIUsesOnlyValuesReturnedByAI(t *testing.T) {
	result := hasilAIFromResponse(7, &dto.AIResponse{CleanedText: "korslet listrik lab", Tokens: []string{"korslet", "listrik"}, Score: 0, Sentimen: "Netral", Urgensi: "Tinggi"})
	if result.CleanedText != "korslet listrik lab" || result.SkorSentimen != 0 || result.Sentimen != "Netral" || result.Urgensi != "Tinggi" {
		t.Fatalf("actual AI fields not preserved: %#v", result)
	}
	if !bytes.Equal(result.Tokens, []byte(`["korslet","listrik"]`)) {
		t.Fatalf("tokens not preserved: %s", result.Tokens)
	}
	if !bytes.Equal(result.DetailSkor, []byte("[]")) {
		t.Fatalf("detail skor must be empty when AI sends no weights: %s", result.DetailSkor)
	}
	if result.SkorPositif == nil || result.SkorNegatif == nil || *result.SkorPositif+*result.SkorNegatif != result.SkorSentimen {
		t.Fatal("sign decomposition must preserve the actual aggregate score")
	}
	if result.PenjelasanSentimen == "" {
		t.Fatal("sentiment explanation must be stored")
	}
}
