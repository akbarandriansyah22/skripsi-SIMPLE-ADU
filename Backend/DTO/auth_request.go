package dto

type RegisterRequest struct {
	NamaLengkap  string `json:"nama_lengkap" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`

	NIM          string `json:"nim" binding:"required"`
	ProgramStudi string `json:"program_studi" binding:"required"`
	Angkatan     int    `json:"angkatan" binding:"required"`
	NoHP         string `json:"no_hp"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}