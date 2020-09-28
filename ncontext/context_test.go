package ncontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextOthers(t *testing.T) {
	m := NewMDC()
	v := m.Other("notexist")
	assert.NotNil(t, v)
	t.Log(v)
}
