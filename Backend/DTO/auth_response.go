package dto

type UserResponse struct {
	ID          uint   `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	IsActive    bool   `json:"is_active"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
