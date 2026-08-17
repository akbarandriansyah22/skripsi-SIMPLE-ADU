package services

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/utils"
	"gorm.io/gorm"
)

var importHeaders = []string{"nama_lengkap", "nim", "email", "program_studi", "angkatan"}

type importStudentRow struct{ Values map[string]string }

func normalizeImportHeader(value string) string {
	value = strings.TrimPrefix(strings.TrimSpace(value), "\ufeff")
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer(" ", "_", "-", "_", ".", "_").Replace(value)
	switch value {
	case "nama", "nama_mahasiswa":
		return "nama_lengkap"
	case "prodi":
		return "program_studi"
	default:
		return value
	}
}

func parseImportRows(reader io.Reader, extension string) ([]importStudentRow, error) {
	extension = strings.ToLower(filepath.Ext(extension))
	if extension == ".csv" {
		csvReader := csv.NewReader(reader)
		csvReader.TrimLeadingSpace = true
		rows, err := csvReader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("CSV tidak valid: %w", err)
		}
		return rowsToStudents(rows)
	}
	if extension != ".xlsx" {
		return nil, errors.New("format file harus CSV atau XLSX")
	}
	data, err := io.ReadAll(io.LimitReader(reader, 12<<20))
	if err != nil {
		return nil, err
	}
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, errors.New("XLSX tidak valid")
	}
	var sheet []byte
	var shared []string
	for _, file := range archive.File {
		if file.Name == "xl/sharedStrings.xml" {
			shared, _ = readSharedStrings(file)
		}
		if file.Name == "xl/worksheets/sheet1.xml" {
			opened, openErr := file.Open()
			if openErr != nil {
				return nil, openErr
			}
			sheet, err = io.ReadAll(opened)
			opened.Close()
			if err != nil {
				return nil, err
			}
		}
	}
	if len(sheet) == 0 {
		return nil, errors.New("sheet pertama XLSX tidak ditemukan")
	}
	return rowsToStudents(parseSheetRows(sheet, shared))
}

func rowsToStudents(rows [][]string) ([]importStudentRow, error) {
	if len(rows) < 2 {
		return nil, errors.New("file import harus memiliki header dan minimal satu baris")
	}
	positions := map[string]int{}
	for i, value := range rows[0] {
		positions[normalizeImportHeader(value)] = i
	}
	for _, required := range importHeaders {
		if _, ok := positions[required]; !ok {
			return nil, fmt.Errorf("kolom wajib %q tidak ditemukan", required)
		}
	}
	result := make([]importStudentRow, 0, len(rows)-1)
	for _, values := range rows[1:] {
		row := importStudentRow{Values: map[string]string{}}
		blank := true
		for _, required := range importHeaders {
			index := positions[required]
			if index < len(values) {
				row.Values[required] = strings.TrimSpace(values[index])
				if row.Values[required] != "" {
					blank = false
				}
			}
		}
		if !blank {
			result = append(result, row)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("tidak ada baris data yang dapat diimpor")
	}
	return result, nil
}

type xlsxCell struct {
	Ref    string `xml:"r,attr"`
	Type   string `xml:"t,attr"`
	Value  string `xml:"v"`
	Inline string `xml:"is>t"`
}
type xlsxRow struct {
	Cells []xlsxCell `xml:"c"`
}
type xlsxSheet struct {
	Rows []xlsxRow `xml:"sheetData>row"`
}
type xlsxShared struct {
	Values []string `xml:"si>t"`
}

func readSharedStrings(file *zip.File) ([]string, error) {
	opened, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer opened.Close()
	var data xlsxShared
	if err := xml.NewDecoder(opened).Decode(&data); err != nil {
		return nil, err
	}
	return data.Values, nil
}
func parseSheetRows(data []byte, shared []string) [][]string {
	var sheet xlsxSheet
	if xml.Unmarshal(data, &sheet) != nil {
		return nil
	}
	rows := make([][]string, 0, len(sheet.Rows))
	for _, row := range sheet.Rows {
		values := make([]string, 0, len(row.Cells))
		for _, cell := range row.Cells {
			value := cell.Value
			if cell.Type == "s" {
				index, _ := strconv.Atoi(value)
				if index >= 0 && index < len(shared) {
					value = shared[index]
				}
			}
			if cell.Inline != "" {
				value = cell.Inline
			}
			col := excelColumnIndex(cell.Ref)
			for len(values) < col {
				values = append(values, "")
			}
			if col < len(values) {
				values[col] = value
			} else {
				values = append(values, value)
			}
		}
		rows = append(rows, values)
	}
	return rows
}

func excelColumnIndex(ref string) int {
	letters := strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return -1
	}, ref)
	if letters == "" {
		return 0
	}
	index := 0
	for _, r := range letters {
		index = index*26 + int(r-'A') + 1
	}
	if index == 0 {
		return 0
	}
	return index - 1
}

func temporaryPassword() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"
	result := make([]byte, 12)
	for index := range result {
		value, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		result[index] = alphabet[value.Int64()]
	}
	return string(result), nil
}

func (s *AdminSistemService) ImportMahasiswa(importedBy uint, filename, extension string, reader io.Reader) (*dto.ImportMahasiswaResponse, error) {
	rows, err := parseImportRows(reader, extension)
	if err != nil {
		return nil, err
	}
	result := &dto.ImportMahasiswaResponse{TotalRows: len(rows), Rows: make([]dto.ImportMahasiswaRowResponse, 0, len(rows))}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		batch := &models.ImportMahasiswa{
			AdminSistemID: importedBy,
			NamaFile:      filepath.Base(filename),
			Status:        "Diproses",
			TotalData:     len(rows),
		}
		if err := tx.Create(batch).Error; err != nil {
			return err
		}
		result.BatchID = batch.ID
		for index, row := range rows {
			item := dto.ImportMahasiswaRowResponse{RowNumber: index + 2, NIM: row.Values["nim"], Email: strings.ToLower(row.Values["email"])}
			savepoint := fmt.Sprintf("sp_row_%d", item.RowNumber)
			if err := tx.SavePoint(savepoint).Error; err != nil {
				return err
			}
			password, userID, rowErr := importOneStudent(tx, row.Values)
			detail := detailFromImportRow(batch.ID, item.RowNumber, row.Values)
			if rowErr != nil {
				if err := tx.RollbackTo(savepoint).Error; err != nil {
					return err
				}
				item.Status = "gagal"
				item.Reason = rowErr.Error()
				result.FailedRows++
				detail.Status = "Gagal"
				detail.UserID = nil
				detail.PesanError = optionalString(item.Reason)
			} else {
				item.Status = "berhasil"
				item.TemporaryPassword = password
				result.SuccessRows++
				detail.Status = "Berhasil"
				detail.UserID = &userID
				detail.PesanError = nil
			}
			result.Rows = append(result.Rows, item)
			if err := tx.Create(detail).Error; err != nil {
				return err
			}
		}
		now := time.Now()
		return tx.Model(batch).Updates(map[string]any{
			"status":          "Selesai",
			"completed_at":    &now,
			"total_data":      result.TotalRows,
			"jumlah_berhasil": result.SuccessRows,
			"jumlah_gagal":    result.FailedRows,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AdminSistemService) ImportHistory() ([]dto.ImportMahasiswaBatchResponse, error) {
	var batches []models.ImportMahasiswa
	if err := config.DB.Order("created_at DESC").Find(&batches).Error; err != nil {
		return nil, err
	}
	result := make([]dto.ImportMahasiswaBatchResponse, 0, len(batches))
	for _, batch := range batches {
		result = append(result, mapImportBatch(batch, nil))
	}
	return result, nil
}

func (s *AdminSistemService) ImportHistoryDetail(id uint64) (*dto.ImportMahasiswaBatchResponse, error) {
	var batch models.ImportMahasiswa
	if err := config.DB.First(&batch, id).Error; err != nil {
		return nil, err
	}
	var rows []models.DetailImportMahasiswa
	if err := config.DB.Where("import_id = ?", batch.ID).Order("nomor_baris ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	mappedRows := make([]dto.ImportMahasiswaRowResponse, 0, len(rows))
	for _, row := range rows {
		mappedRows = append(mappedRows, dto.ImportMahasiswaRowResponse{
			RowNumber: row.NomorBaris,
			NIM:       stringValue(row.NIM),
			Email:     stringValue(row.Email),
			Status:    strings.ToLower(row.Status),
			Reason:    stringValue(row.PesanError),
		})
	}
	result := mapImportBatch(batch, mappedRows)
	return &result, nil
}

func mapImportBatch(batch models.ImportMahasiswa, rows []dto.ImportMahasiswaRowResponse) dto.ImportMahasiswaBatchResponse {
	return dto.ImportMahasiswaBatchResponse{
		ID:          batch.ID,
		ImportedBy:  batch.AdminSistemID,
		FileName:    batch.NamaFile,
		TotalRows:   batch.TotalData,
		SuccessRows: batch.JumlahBerhasil,
		FailedRows:  batch.JumlahGagal,
		CreatedAt:   batch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Rows:        rows,
	}
}

func detailFromImportRow(importID uint, nomorBaris int, values map[string]string) *models.DetailImportMahasiswa {
	detail := &models.DetailImportMahasiswa{
		ImportID:     importID,
		NomorBaris:   nomorBaris,
		NamaLengkap:  optionalString(values["nama_lengkap"]),
		NIM:          optionalString(values["nim"]),
		Email:        optionalString(strings.ToLower(values["email"])),
		ProgramStudi: optionalString(values["program_studi"]),
	}
	if angkatan, err := strconv.Atoi(strings.TrimSpace(values["angkatan"])); err == nil && angkatan >= 2000 && angkatan <= 2100 {
		detail.Angkatan = &angkatan
	}
	return detail
}

func importOneStudent(tx *gorm.DB, values map[string]string) (string, uint, error) {
	for _, field := range importHeaders {
		if strings.TrimSpace(values[field]) == "" {
			return "", 0, fmt.Errorf("kolom %s wajib diisi", field)
		}
	}
	if !strings.Contains(values["email"], "@") {
		return "", 0, errors.New("format email tidak valid")
	}
	angkatan, err := strconv.Atoi(values["angkatan"])
	if err != nil || angkatan < 2000 || angkatan > 2100 {
		return "", 0, errors.New("angkatan harus berupa tahun yang valid")
	}
	var count int64
	if err := tx.Model(&models.User{}).Where("LOWER(email) = ?", strings.ToLower(values["email"])).Count(&count).Error; err != nil {
		return "", 0, err
	}
	if count > 0 {
		return "", 0, errors.New("email sudah terdaftar")
	}
	if err := tx.Model(&models.Mahasiswa{}).Where("nim = ?", values["nim"]).Count(&count).Error; err != nil {
		return "", 0, err
	}
	if count > 0 {
		return "", 0, errors.New("NIM sudah terdaftar")
	}
	password, err := temporaryPassword()
	if err != nil {
		return "", 0, err
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", 0, err
	}
	user := &models.User{
		NamaLengkap:        values["nama_lengkap"],
		Email:              strings.ToLower(values["email"]),
		PasswordHash:       hash,
		Role:               utils.RoleMahasiswa,
		IsActive:           true,
		PasswordMustChange: true,
		SumberAkun:         "import",
	}
	if err := tx.Create(user).Error; err != nil {
		return "", 0, err
	}
	student := &models.Mahasiswa{UserID: user.ID, NIM: values["nim"], ProgramStudi: values["program_studi"], Angkatan: angkatan}
	if err := tx.Create(student).Error; err != nil {
		return "", 0, err
	}
	return password, user.ID, nil
}
