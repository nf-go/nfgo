package nconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptSecretValue(t *testing.T) {
	a := assert.New(t)
	secretKey = "Uv38ByGCZU8WP18PmmIdcpVmx00QA3xNe7sEB9Hixkk="
	encryptedText := "SECRET(gYVa2GgdDYbR6R4AlvB2Tcwi/z9yT20tkiADAsd4yxO3flA1xV0a)"

	plainText, err := DecryptSecretValue(encryptedText)
	a.Nil(err)
	a.Equal("hello world", plainText)

	encryptedText = "gYVa2GgdDYbR6R4AlvB2Tcwi/z9yT20tkiADAsd4yxO3flA1xV0a"
	plainText, err = DecryptSecretValue(encryptedText)
	a.Nil(err)
	a.Equal("gYVa2GgdDYbR6R4AlvB2Tcwi/z9yT20tkiADAsd4yxO3flA1xV0a", plainText)
}
