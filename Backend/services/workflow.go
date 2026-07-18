package services

import "strings"

const (
	StatusMenunggu   = "Menunggu Verifikasi"
	StatusDitolak    = "Ditolak"
	StatusDiteruskan = "Diteruskan ke Pimpinan"
	StatusDiproses   = "Diproses"
	StatusSelesai    = "Selesai"
)

func isValidStatusTransition(current, next string) bool {
	if current == next {
		return false
	}
	switch strings.ToLower(current) {
	case strings.ToLower(StatusMenunggu):
		return strings.EqualFold(next, StatusDitolak) || strings.EqualFold(next, StatusDiteruskan) || strings.EqualFold(next, StatusDiproses)
	case strings.ToLower(StatusDiteruskan):
		return strings.EqualFold(next, StatusDiproses)
	case strings.ToLower(StatusDiproses):
		return strings.EqualFold(next, StatusSelesai)
	default:
		return false
	}
}
