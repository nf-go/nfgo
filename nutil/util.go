package nutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"nfgo.ga/nfgo/nutil/id"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Sha256 -
func Sha256(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

// RandString -
func RandString(length int) string {
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

// UUID -
func UUID() (string, error) {
	uuid, err := id.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.Hex(), err
}

// RedisKey -
type RedisKey string

// Key -
func (k RedisKey) Key(a ...interface{}) string {
	return fmt.Sprintf(string(k), a...)
}
