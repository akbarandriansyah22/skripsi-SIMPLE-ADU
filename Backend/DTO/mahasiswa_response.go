package dto

type ProfileMahasiswaResponse struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	NIM           string `json:"nim"`
	ProgramStudi  string `json:"program_studi"`
	Angkatan      int    `json:"angkatan"`
	NoHP          string `json:"no_hp"`

	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	IsActive    bool   `json:"is_active"`
}