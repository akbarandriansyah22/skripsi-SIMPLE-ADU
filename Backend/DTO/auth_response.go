package dto

type UserResponse struct {
	ID                 uint                       `json:"id"`
	NamaLengkap        string                     `json:"nama_lengkap"`
	Email              string                     `json:"email"`
	Role               string                     `json:"role"`
	UnitID             *uint                      `json:"unit_id,omitempty"`
	UnitName           string                     `json:"unit_name,omitempty"`
	IsActive           bool                       `json:"is_active"`
	PasswordMustChange bool                       `json:"password_must_change,omitempty"`
	Mahasiswa          *MahasiswaIdentityResponse `json:"mahasiswa,omitempty"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
