package ncontext

import "testing"

func TestContextOthers(t *testing.T) {
	m := NewMDC()
	v := m.Other("notexist")
	if v != nil {
		t.FailNow()
	}
	t.Log(v)
}
