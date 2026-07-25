package dto

type CreatePengaduanRequest struct {
	KategoriID       uint   `json:"kategori_id" binding:"required"`
	Judul            string `json:"judul" binding:"required"`
	Deskripsi        string `json:"deskripsi" binding:"required"`
	Lampiran         string `json:"lampiran"`
	LampiranNamaAsli string `json:"lampiran_nama_asli"`
	LampiranMimeType string `json:"lampiran_mime_type"`
	LampiranUkuran   int64  `json:"lampiran_ukuran"`
}

type UpdatePengaduanRequest struct {
	KategoriID       uint   `json:"kategori_id"`
	Judul            string `json:"judul"`
	Deskripsi        string `json:"deskripsi"`
	Lampiran         string `json:"lampiran"`
	LampiranNamaAsli string `json:"lampiran_nama_asli"`
	LampiranMimeType string `json:"lampiran_mime_type"`
	LampiranUkuran   int64  `json:"lampiran_ukuran"`
}
