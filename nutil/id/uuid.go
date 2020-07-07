package id

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// UUID - represents a UUID.
type UUID [16]byte

// NewUUID - generates a new uuid.
func NewUUID() (UUID, error) {
	var uuid [16]byte

	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return [16]byte{}, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return uuid, nil
}

// Hex - returns a hex representation of the UUID.
func (id UUID) Hex() string {
	return hex.EncodeToString(id[:])
}
