package nutil

import (
	"testing"
)

func TestSha256(t *testing.T) {
	plaint := "qwer1234JSQotTd2s3"
	hashed := Sha256(plaint)
	if hashed != "003b12dfe2c657c572f3496a63af8e2250f465e96aa0ebfea8428c82e0b10ba0" {
		t.FailNow()
	}
	t.Log(len(hashed))
}

func TestRandomString(t *testing.T) {
	str := RandString(6)
	t.Log(str)
	if len(str) != 6 {
		t.FailNow()
	}
}
