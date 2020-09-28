package nutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256(t *testing.T) {
	plaint := "qwer1234JSQotTd2s3"
	hashed := Sha256(plaint)
	assert.Equal(t, "003b12dfe2c657c572f3496a63af8e2250f465e96aa0ebfea8428c82e0b10ba0", hashed)
	t.Log(len(hashed))
}

func TestRandomString(t *testing.T) {
	str := RandString(6)
	t.Log(str)
	assert.Equal(t, 6, len(str))
}

func TestIsNil(t *testing.T) {
	var ch chan int
	if !IsNil(ch) {
		t.Fatal()
	}

	var f func() = nil
	if !IsNil(f) {
		t.Fatal()
	}

	var m map[string]int = nil
	if !IsNil(m) {
		t.Fatal()
	}

	var s []int = nil
	if !IsNil(s) {
		t.Fatal(s)
	}

	var obj interface{} = (*int)(nil)
	if !IsNil(obj) {
		t.Fatal()
	}

	var obj2 interface{} = nil
	if !IsNil(obj2) {
		t.Fatal()
	}

	if IsNil("") || IsNil(1) {
		t.Fatal()
	}

}
