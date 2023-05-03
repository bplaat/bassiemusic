package database

import (
	"database/sql"
	"encoding/json"
)

// Null string
type NullString struct {
	sql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte("null"), nil
}

// Null float
type NullFloat64 struct {
	sql.NullFloat64
}

func (s NullFloat64) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Float64)
	}
	return []byte("null"), nil
}
