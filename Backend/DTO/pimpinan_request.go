package dto

type DisposisiRequest struct {
	Catatan string `json:"catatan" binding:"required"`
}
