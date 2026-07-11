package dto

type DashboardPimpinanResponse struct {
	TotalUrgensiTinggi int64 `json:"total_urgensi_tinggi"`
	BelumDisposisi     int64 `json:"belum_disposisi"`
	SudahDisposisi     int64 `json:"sudah_disposisi"`
}
