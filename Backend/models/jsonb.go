package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONB maps PostgreSQL jsonb columns without allowing GORM to alter the
// already-provisioned schema.
type JSONB json.RawMessage

func (j *JSONB) Scan(value any) error {
	switch typed := value.(type) {
	case nil:
		*j = nil
		return nil
	case []byte:
		*j = append((*j)[:0], typed...)
		return nil
	case string:
		*j = append((*j)[:0], typed...)
		return nil
	default:
		return fmt.Errorf("models.JSONB: unsupported database value %T", value)
	}
}

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	if !json.Valid(j) {
		return nil, fmt.Errorf("models.JSONB: invalid JSON")
	}
	return []byte(j), nil
}

func (JSONB) GormDataType() string { return "jsonb" }
