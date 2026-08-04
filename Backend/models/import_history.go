package models

import "time"

type ImportMahasiswaBatch struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ImportedBy  uint      `gorm:"column:imported_by;not null" json:"imported_by"`
	FileName    string    `gorm:"column:file_name;size:255;not null" json:"file_name"`
	TotalRows   int       `gorm:"column:total_rows;not null" json:"total_rows"`
	SuccessRows int       `gorm:"column:success_rows;not null" json:"success_rows"`
	FailedRows  int       `gorm:"column:failed_rows;not null" json:"failed_rows"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

type ImportMahasiswaRow struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BatchID   uint      `gorm:"column:batch_id;not null;index" json:"batch_id"`
	RowNumber int       `gorm:"column:row_number;not null" json:"row_number"`
	NIM       string    `gorm:"column:nim;size:50" json:"nim"`
	Email     string    `gorm:"column:email;size:150" json:"email"`
	Status    string    `gorm:"column:status;size:20;not null" json:"status"`
	Reason    string    `gorm:"column:reason;type:text" json:"reason"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ImportMahasiswaBatch) TableName() string { return "import_mahasiswa_batch" }
func (ImportMahasiswaRow) TableName() string   { return "import_mahasiswa_row" }
