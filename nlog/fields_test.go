package nlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFields(t *testing.T) {
	fs := NewFields()
	assert.Equal(t, 0, len(fs))

	fs = NewFields("k1")
	assert.Equal(t, 0, len(fs))

	fs = NewFields("k1", "v1")
	assert.Equal(t, 1, len(fs))
	assert.Equal(t, "v1", fs["k1"])

	fs = NewFields("k1", "v1", "k2")
	assert.Equal(t, 1, len(fs))
	assert.Equal(t, "v1", fs["k1"])

	fs = NewFields("k1", "v1", "k2", "v2")
	assert.Equal(t, 2, len(fs))
	assert.Equal(t, "v1", fs["k1"])
	assert.Equal(t, "v2", fs["k2"])

	fs = NewFields("k1", "", "k2", "")
	assert.Equal(t, 0, len(fs))
}
