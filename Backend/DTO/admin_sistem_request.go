package dto

type CreateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        string `json:"role" binding:"required"`
	UnitID      *uint  `json:"unit_id"`
}

type UpdateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Role        string `json:"role" binding:"required"`
	UnitID      *uint  `json:"unit_id"`
}

type UserStatusRequest struct {
	IsActive bool `json:"is_active"`
}
type ActiveStatusRequest struct {
	IsActive bool `json:"is_active"`
}
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

type UnitRequest struct {
	NamaUnit string `json:"nama_unit" binding:"required"`
	Email    string `json:"email"`
}

type CategoryRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

type ImportMahasiswaRowResponse struct {
	RowNumber         int    `json:"row_number"`
	NIM               string `json:"nim,omitempty"`
	Email             string `json:"email,omitempty"`
	Status            string `json:"status"`
	Reason            string `json:"reason,omitempty"`
	TemporaryPassword string `json:"temporary_password,omitempty"`
}

type ImportMahasiswaResponse struct {
	BatchID     uint                         `json:"batch_id"`
	TotalRows   int                          `json:"total_rows"`
	SuccessRows int                          `json:"success_rows"`
	FailedRows  int                          `json:"failed_rows"`
	Rows        []ImportMahasiswaRowResponse `json:"rows"`
}
