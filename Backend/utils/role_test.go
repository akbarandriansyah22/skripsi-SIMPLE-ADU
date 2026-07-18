package utils

import "testing"

func TestCanonicalRole(t *testing.T) {
	tests := map[string]string{
		"mahasiswa":         RoleMahasiswa,
		"petugas":           RoleAdmin,
		"Admin Fakultas":    RoleAdmin,
		"admin_fakultas":    RoleAdmin,
		"pimpinan":          RolePimpinan,
		"Pimpinan Fakultas": RolePimpinan,
		"pimpinan_fakultas": RolePimpinan,
	}
	for input, expected := range tests {
		if actual := CanonicalRole(input); actual != expected {
			t.Fatalf("CanonicalRole(%q) = %q, expected %q", input, actual, expected)
		}
	}
	if IsCanonicalRole("unknown") {
		t.Fatal("unknown role must not be accepted")
	}
}
