package dto

type DashboardAdminResponse struct {
	TotalPengaduan int64 `json:"total_pengaduan"`
	Menunggu       int64 `json:"menunggu"`
	Diproses       int64 `json:"diproses"`
	Selesai        int64 `json:"selesai"`
	Ditolak        int64 `json:"ditolak"`
}
