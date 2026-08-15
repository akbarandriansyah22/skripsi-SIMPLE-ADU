package services

import "testing"

func TestComplaintStatusTransitionMatrix(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{StatusMenunggu, StatusDitolak, true},
		{StatusMenunggu, StatusMenungguDisposisi, true},
		{StatusMenunggu, StatusDiteruskanUnit, true},
		{StatusMenungguDisposisi, StatusDiteruskanUnit, true},
		{StatusMenungguDisposisi, StatusDitolak, true},
		{StatusDiteruskanUnit, StatusDiproses, true},
		{StatusDiteruskanUnit, StatusMenunggu, true},
		{StatusDiproses, StatusSelesai, true},
		{StatusDiproses, StatusMenunggu, true},
		{StatusDitolak, StatusDiproses, false},
		{StatusSelesai, StatusDiproses, false},
		{StatusMenungguDisposisi, StatusDiproses, false},
	}
	for _, test := range tests {
		if got := isValidStatusTransition(test.from, test.to); got != test.want {
			t.Errorf("transition %q -> %q = %v, want %v", test.from, test.to, got, test.want)
		}
	}
}
