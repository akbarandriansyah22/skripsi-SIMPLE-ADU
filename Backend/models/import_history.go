package models

import "time"

type ImportMahasiswa struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminSistemID  uint       `gorm:"column:admin_sistem_id;not null" json:"admin_sistem_id"`
	NamaFile       string     `gorm:"column:nama_file;size:255;not null" json:"nama_file"`
	Status         string     `gorm:"column:status;size:20;not null;default:Diproses" json:"status"`
	TotalData      int        `gorm:"column:total_data;not null" json:"total_data"`
	JumlahBerhasil int        `gorm:"column:jumlah_berhasil;not null" json:"jumlah_berhasil"`
	JumlahGagal    int        `gorm:"column:jumlah_gagal;not null" json:"jumlah_gagal"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
}

type DetailImportMahasiswa struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ImportID     uint      `gorm:"column:import_id;not null;index" json:"import_id"`
	NomorBaris   int       `gorm:"column:nomor_baris;not null" json:"nomor_baris"`
	NamaLengkap  *string   `gorm:"column:nama_lengkap;size:150" json:"nama_lengkap,omitempty"`
	NIM          *string   `gorm:"column:nim;size:20" json:"nim,omitempty"`
	Email        *string   `gorm:"column:email;size:150" json:"email,omitempty"`
	ProgramStudi *string   `gorm:"column:program_studi;size:100" json:"program_studi,omitempty"`
	Angkatan     *int      `gorm:"column:angkatan" json:"angkatan,omitempty"`
	Status       string    `gorm:"column:status;size:20;not null" json:"status"`
	UserID       *uint     `gorm:"column:user_id" json:"user_id,omitempty"`
	PesanError   *string   `gorm:"column:pesan_error;type:text" json:"pesan_error,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ImportMahasiswa) TableName() string       { return "import_mahasiswa" }
func (DetailImportMahasiswa) TableName() string { return "detail_import_mahasiswa" }
