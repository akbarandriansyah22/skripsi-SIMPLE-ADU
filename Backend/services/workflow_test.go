package services

import "testing"

func TestComplaintStatusTransitionMatrix(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		{StatusMenunggu, StatusDitolak, true},
		{StatusMenunggu, StatusDiteruskan, true},
		{StatusMenunggu, StatusDiproses, true},
		{StatusDiteruskan, StatusDiproses, true},
		{StatusDiproses, StatusSelesai, true},
		{StatusDitolak, StatusDiproses, false},
		{StatusSelesai, StatusDiproses, false},
		{StatusDiteruskan, StatusSelesai, false},
	}
	for _, test := range tests {
		if got := isValidStatusTransition(test.from, test.to); got != test.want {
			t.Errorf("transition %q -> %q = %v, want %v", test.from, test.to, got, test.want)
		}
	}
}
