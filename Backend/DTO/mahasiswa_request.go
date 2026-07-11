package dto

type UpdateProfileRequest struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	NoHP        string `json:"no_hp" binding:"required"`
}
