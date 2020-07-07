package jwt

import (
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	jwtSecret := "VoCra8#GEBAbRl*+vos9UF@??gi8Oy"
	token, err := NewToken(jwtSecret, "admin", time.Now().Add(time.Hour*24*365*20), map[string]interface{}{"userId": 1})
	t.Log(token)
	if err != nil {
		t.FailNow()
	}
}

func TestParseToken(t *testing.T) {
	jwtSecret := "VoCra8#GEBAbRl*+vos9UF@??gi8Oy"
	token := "eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjIyMjY3MTIzMDkuODY5ODM5Nywic3ViIjoiYWRtaW4iLCJ1c2VySWQiOjF9.2M1tC3oQrC0Ym7K9qfVVcLS-eyCM4bacVOyuh6Jj1yg"
	claims, err := ParseToken(jwtSecret, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims.Subject)
	t.Log(claims.Expires)
	t.Log(claims.Set)
}
