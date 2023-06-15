package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Uuid
type Uuid struct {
	Bytes [16]byte
}

func New() Uuid {
	bytes := [16]byte{}
	if _, err := rand.Read(bytes[:]); err != nil {
		log.Fatalln(err)
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return Uuid{bytes}
}

func IsValid(uuid string) bool {
	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return re.MatchString(uuid)
}

func Parse(uuid string) (Uuid, error) {
	if !IsValid(uuid) {
		return Uuid{}, fmt.Errorf("invalid uuid")
	}

	bytes, err := hex.DecodeString(strings.ReplaceAll(uuid, "-", ""))
	if err != nil {
		log.Fatalln(err)
	}
	parsedUuid := Uuid{}
	copy(parsedUuid.Bytes[:], bytes)
	return parsedUuid, nil
}

func (uuid Uuid) Equals(rhs Uuid) bool {
	for i := 0; i < 16; i++ {
		if uuid.Bytes[i] != rhs.Bytes[i] {
			return false
		}
	}
	return true
}

func (uuid *Uuid) Scan(value any) error {
	copy(uuid.Bytes[:], value.([]byte))
	return nil
}

func (uuid Uuid) Value() (driver.Value, error) {
	return uuid.Bytes[:], nil
}

func (uuid Uuid) String() string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		uuid.Bytes[0], uuid.Bytes[1], uuid.Bytes[2], uuid.Bytes[3], uuid.Bytes[4], uuid.Bytes[5], uuid.Bytes[6], uuid.Bytes[7],
		uuid.Bytes[8], uuid.Bytes[9], uuid.Bytes[10], uuid.Bytes[11], uuid.Bytes[12], uuid.Bytes[13], uuid.Bytes[14], uuid.Bytes[15])
}

func (uuid Uuid) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.String())
}

// NullUuid
type NullUuid struct {
	Uuid  Uuid
	Valid bool
}

func (uuid *NullUuid) Scan(value any) error {
	if value != nil {
		uuid.Uuid.Scan(value)
		uuid.Valid = true
	}
	return nil
}

func (uuid NullUuid) Value() (driver.Value, error) {
	if uuid.Valid {
		return uuid.Uuid.Value()
	}
	return nil, nil
}

func (uuid NullUuid) String() string {
	if uuid.Valid {
		return uuid.Uuid.String()
	}
	return "nil"
}

func (uuid NullUuid) MarshalJSON() ([]byte, error) {
	if uuid.Valid {
		return json.Marshal(uuid.Uuid)
	}
	return []byte("null"), nil
}
