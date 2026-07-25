package models

import "time"

type HasilAI struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID uint   `gorm:"column:pengaduan_id;unique;not null" json:"pengaduan_id"`
	CleanedText string `gorm:"column:cleaned_text;type:text" json:"cleaned_text"`
	Tokens      JSONB  `gorm:"column:tokens;type:jsonb" json:"tokens"`
	// The final schema requires these columns and supplies database default 0.
	// A zero here means the AI response did not provide a separately computable
	// positive/negative breakdown; it is never presented as a calculated score.
	SkorPositif        *int      `gorm:"column:skor_positif;default:0" json:"skor_positif,omitempty"`
	SkorNegatif        *int      `gorm:"column:skor_negatif;default:0" json:"skor_negatif,omitempty"`
	SkorSentimen       int       `gorm:"column:skor_sentimen;not null" json:"skor_sentimen"`
	Sentimen           string    `gorm:"column:sentimen;size:20;not null" json:"sentimen"`
	PenjelasanSentimen string    `gorm:"column:penjelasan_sentimen;type:text" json:"penjelasan_sentimen"`
	DetailSkor         JSONB     `gorm:"column:detail_skor;type:jsonb" json:"detail_skor"`
	Urgensi            string    `gorm:"column:urgensi;size:20;not null" json:"urgensi"`
	DasarUrgensi       string    `gorm:"column:dasar_urgensi;type:text" json:"dasar_urgensi"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relasi
	Pengaduan *Pengaduan `gorm:"foreignKey:PengaduanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (HasilAI) TableName() string {
	return "hasil_ai"
}
