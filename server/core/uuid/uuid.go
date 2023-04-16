package uuid

import (
	"crypto/rand"
	"fmt"
	"log"
	"regexp"
)

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

func (uuid Uuid) String() string {
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		uuid.Bytes[0], uuid.Bytes[1], uuid.Bytes[2], uuid.Bytes[3], uuid.Bytes[4], uuid.Bytes[5], uuid.Bytes[6], uuid.Bytes[7],
		uuid.Bytes[8], uuid.Bytes[9], uuid.Bytes[10], uuid.Bytes[11], uuid.Bytes[12], uuid.Bytes[13], uuid.Bytes[14], uuid.Bytes[15])
}

func IsValid(uuid string) bool {
	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return re.MatchString(uuid)
}
