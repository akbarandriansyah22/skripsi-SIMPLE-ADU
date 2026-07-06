package dto

type CreatePengaduanRequest struct {
	KategoriID uint   `json:"kategori_id" binding:"required"`
	Judul      string `json:"judul" binding:"required"`
	Deskripsi  string `json:"deskripsi" binding:"required"`
	Lampiran   string `json:"lampiran"`
}

type UpdatePengaduanRequest struct {
	KategoriID uint   `json:"kategori_id"`
	Judul      string `json:"judul"`
	Deskripsi  string `json:"deskripsi"`
	Lampiran   string `json:"lampiran"`
}