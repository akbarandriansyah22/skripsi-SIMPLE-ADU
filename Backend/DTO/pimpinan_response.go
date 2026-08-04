package dto

type DashboardPimpinanResponse struct {
	TotalPengaduan     int64 `json:"total_pengaduan"`
	BelumDikerjakan    int64 `json:"belum_dikerjakan"`
	SedangDiproses     int64 `json:"sedang_diproses"`
	Selesai            int64 `json:"selesai"`
	Ditolak            int64 `json:"ditolak"`
	UrgensiRendah      int64 `json:"urgensi_rendah"`
	UrgensiSedang      int64 `json:"urgensi_sedang"`
	TotalUrgensiTinggi int64 `json:"total_urgensi_tinggi"`
	BelumDisposisi     int64 `json:"belum_disposisi"`
	SudahDisposisi     int64 `json:"sudah_disposisi"`
}
