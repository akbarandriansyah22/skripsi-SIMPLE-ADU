package utils

import "strings"

const (
	RoleMahasiswa = "mahasiswa"
	RoleAdmin     = "admin_fakultas"
	RolePimpinan  = "pimpinan_fakultas"
)

// CanonicalRole keeps persisted roles stable while accepting roles written by
// earlier versions of the application.
func CanonicalRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "mahasiswa":
		return RoleMahasiswa
	case "petugas", "admin", "admin fakultas", "admin_fakultas":
		return RoleAdmin
	case "pimpinan", "pimpinan fakultas", "pimpinan_fakultas":
		return RolePimpinan
	default:
		return ""
	}
}

func IsCanonicalRole(role string) bool {
	return CanonicalRole(role) != ""
}
