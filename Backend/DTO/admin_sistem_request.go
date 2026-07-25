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
