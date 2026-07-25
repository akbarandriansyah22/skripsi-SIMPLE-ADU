package dto

type DisposisiRequest struct {
	UnitID  uint   `json:"unit_id" binding:"required"`
	Catatan string `json:"catatan" binding:"required"`
}
