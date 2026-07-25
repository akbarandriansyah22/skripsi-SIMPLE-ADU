package dto

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AssignUnitRequest struct {
	UnitID uint `json:"unit_id" binding:"required"`
}

type ValidatePengaduanRequest struct {
	StatusValidasi string `json:"status_validasi" binding:"required"`
	Catatan        string `json:"catatan"`
}

type ForwardDisposisiRequest struct {
	PengaduanID uint   `json:"pengaduan_id" binding:"required"`
	PimpinanID  uint   `json:"pimpinan_id" binding:"required"`
	Catatan     string `json:"catatan"`
}
