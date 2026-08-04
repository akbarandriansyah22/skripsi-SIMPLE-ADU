package dto

import (
	"encoding/json"
	"time"
)

type PengaduanResponse struct {
	ID                   uint                      `json:"id"`
	KodeTiket            string                    `json:"kode_tiket"`
	UserID               uint                      `json:"user_id"`
	KategoriID           uint                      `json:"kategori_id"`
	Kategori             *KategoriResponse         `json:"kategori,omitempty"`
	UnitID               *uint                     `json:"unit_id,omitempty"`
	Unit                 *UnitResponse             `json:"unit,omitempty"`
	User                 *UserResponse             `json:"user,omitempty"`
	Judul                string                    `json:"judul"`
	Deskripsi            string                    `json:"deskripsi"`
	Lampiran             string                    `json:"lampiran"`
	LampiranNamaAsli     string                    `json:"lampiran_nama_asli,omitempty"`
	LampiranMimeType     string                    `json:"lampiran_mime_type,omitempty"`
	LampiranUkuran       int64                     `json:"lampiran_ukuran,omitempty"`
	LampiranURL          string                    `json:"lampiran_url,omitempty"`
	Status               string                    `json:"status"`
	SkorSentimen         *int                      `json:"skor_sentimen,omitempty"`
	Sentimen             string                    `json:"sentimen,omitempty"`
	Urgensi              string                    `json:"urgensi,omitempty"`
	SkorPositif          *int                      `json:"skor_positif,omitempty"`
	SkorNegatif          *int                      `json:"skor_negatif,omitempty"`
	SentimentScore       *int                      `json:"sentiment_score,omitempty"`
	MatchedWords         json.RawMessage           `json:"matched_words,omitempty"`
	CleanedText          string                    `json:"cleaned_text,omitempty"`
	SentimentExplanation string                    `json:"sentiment_explanation,omitempty"`
	UrgencyScore         int                       `json:"urgency_score,omitempty"`
	UrgencyReason        string                    `json:"urgency_reason,omitempty"`
	AIStatus             string                    `json:"ai_status"`
	Validasi             *ValidasiResponse         `json:"validasi,omitempty"`
	Disposisi            *DisposisiResponse        `json:"disposisi,omitempty"`
	RiwayatStatus        []RiwayatStatusResponse   `json:"riwayat_status_pengaduan,omitempty"`
	ResponPengaduan      []ResponPengaduanResponse `json:"respon_pengaduan,omitempty"`
	CreatedAt            time.Time                 `json:"created_at"`
	UpdatedAt            time.Time                 `json:"updated_at"`
	TanggalMasukPimpinan *time.Time                `json:"tanggal_masuk_pimpinan,omitempty"`
	KoordinasiInternal   []KoordinasiResponse      `json:"koordinasi_internal,omitempty"`
}

type KategoriResponse struct {
	ID   uint   `json:"id"`
	Nama string `json:"nama"`
}

type UnitResponse struct {
	ID       uint   `json:"id"`
	NamaUnit string `json:"nama_unit"`
}

type MahasiswaIdentityResponse struct {
	NamaLengkap  string `json:"nama_lengkap"`
	NIM          string `json:"nim"`
	ProgramStudi string `json:"program_studi"`
	Angkatan     int    `json:"angkatan"`
}

type ValidasiResponse struct {
	ID              uint   `json:"id"`
	AdminFakultasID uint   `json:"admin_fakultas_id"`
	StatusValidasi  string `json:"status_validasi"`
	Catatan         string `json:"catatan"`
}

type DisposisiResponse struct {
	ID         uint          `json:"id"`
	PimpinanID uint          `json:"pimpinan_id"`
	UnitID     uint          `json:"unit_id"`
	Catatan    string        `json:"catatan"`
	Unit       *UnitResponse `json:"unit,omitempty"`
}

type RiwayatStatusResponse struct {
	ID         uint      `json:"id"`
	ChangedBy  uint      `json:"changed_by"`
	StatusLama string    `json:"status_lama"`
	StatusBaru string    `json:"status_baru"`
	Catatan    string    `json:"catatan"`
	CreatedAt  time.Time `json:"created_at"`
}

type ReanalyzeAIResponse struct {
	PengaduanID  uint   `json:"pengaduan_id"`
	SkorSentimen int    `json:"skor_sentimen"`
	Sentimen     string `json:"sentimen"`
	Urgensi      string `json:"urgensi"`
	AIStatus     string `json:"ai_status"`
}

type ResponPengaduanResponse struct {
	ID               uint          `json:"id"`
	PengaduanID      uint          `json:"pengaduan_id"`
	UserID           uint          `json:"user_id"`
	Pesan            string        `json:"pesan"`
	Lampiran         string        `json:"lampiran,omitempty"`
	LampiranNamaAsli string        `json:"lampiran_nama_asli,omitempty"`
	LampiranMimeType string        `json:"lampiran_mime_type,omitempty"`
	LampiranUkuran   int64         `json:"lampiran_ukuran,omitempty"`
	LampiranURL      string        `json:"lampiran_url,omitempty"`
	User             *UserResponse `json:"user,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
}
