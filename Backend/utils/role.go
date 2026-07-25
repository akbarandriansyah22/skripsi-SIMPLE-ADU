package utils

import "strings"

const (
	RoleMahasiswa     = "mahasiswa"
	RoleAdminSistem   = "admin_sistem"
	RoleAdminFakultas = "admin_fakultas"
	RoleKasubag       = "kasubag"
	RolePimpinan      = "pimpinan_fakultas"
	// RoleAdmin is kept as a source-compatible alias for existing callers.
	RoleAdmin = RoleAdminFakultas
)

// CanonicalRole keeps persisted roles stable while accepting roles written by
// earlier versions of the application.
func CanonicalRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "mahasiswa":
		return RoleMahasiswa
	case "admin sistem", "admin_sistem":
		return RoleAdminSistem
	case "petugas", "admin", "admin fakultas", "admin_fakultas":
		return RoleAdminFakultas
	case "kasubag", "kasubag akademik", "kasubag sarpras", "kasubag satpras", "kasubag sarana dan prasarana":
		return RoleKasubag
	case "pimpinan", "pimpinan fakultas", "pimpinan_fakultas":
		return RolePimpinan
	default:
		return ""
	}
}

func IsCanonicalRole(role string) bool {
	return CanonicalRole(role) != ""
}
